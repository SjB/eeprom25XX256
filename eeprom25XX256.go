// Package eeprom25XX256 allows to read and write to a 256K SPI Bus Serial EEPROM
package eeprom25XX256

import (
	"encoding/binary"
	"errors"
	"os"
	"time"

	"github.com/kidoman/embd"
)

const (
	twc              = 5 * time.Millisecond //duration of write cycle
	pageSize         = 64
	maxReadBlockSize = 4093 // linux SPIDEV bufsize minus 3 bytes chip overhead
	maxAddress       = 32768
)

const (
	instwrite = 0x02
	instread  = 0x03
	instwrdi  = 0x04
	instwren  = 0x06
	instrdsr  = 0x05
)

// eeprom25XX256 represents a 256K SPI Bus Serial EEPROM
type eeprom25XX256 struct {
	offset uint16
	bus    embd.SPIBus
}

// New creates a new eeprom25xx256 interface.
// The SPIBUS need to be fully configured
func New(bus embd.SPIBus) *eeprom25XX256 {
	return &eeprom25XX256{0, bus}
}

// Seek sets the offset for the next Read or Write on the eeprom, interpreted
// according to whence: 0 means relative to origin of the file, 1 means
// relative to the current offset, and 2 means relative to the end. It returns
// the new offset and an error, if any
func (s *eeprom25XX256) Seek(offset int64, whence int) (ret int64, err error) {
	curOffset := int64(s.offset)
	switch whence {
	case 0:
		curOffset = offset
	case 1:
		curOffset += offset
	case 2:
		curOffset = maxAddress - offset
	}
	if (curOffset < maxAddress) && (curOffset > 0) {
		s.offset = uint16(curOffset)
		return curOffset, nil
	}
	return int64(s.offset), errors.New("Invalid Offset")
}

// Read reads up to len(b) bytes from the eeprom, It returns the number of
// bytes read and an error, if any. The read will loop to the 1 addressable
// location, if the last addressable location is exceeded.
func (s *eeprom25XX256) Read(b []byte) (n int, err error) {

	byteCount := len(b)
	n = 0
	for n < byteCount {
		blockSize := byteCount - n
		if blockSize > maxReadBlockSize {
			blockSize = maxReadBlockSize
		}
		c, err := s.readBlock(b[n : n+blockSize])
		n += c
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (s *eeprom25XX256) readBlock(b []byte) (n int, err error) {

	bufsize := len(b) + 3 // instread + addr
	buf := make([]uint8, bufsize)

	buf[0] = instread
	binary.BigEndian.PutUint16(buf[1:], s.offset)
	if err = s.bus.TransferAndRecieveData(buf[:]); err != nil {
		return 0, err
	}
	copy(b, buf[3:])
	n = len(b)
	s.offset += uint16(n % maxAddress)
	return n, nil
}

// ReadAt read len(b) bytes from the eeprom starting at byte offset off. It
// returns the number of bytes read and the error, if any. The read will loop to
// the 1 addressable location, if the last addressable location is exceeded.
func (s *eeprom25XX256) ReadAt(b []byte, off int64) (n int, err error) {
	if _, err = s.Seek(off, os.SEEK_SET); err != nil {
		return 0, err
	}

	return s.Read(b)
}

// Write writes len(b) bytes to the eeprom. It returns the number of bytes
// written and an error, if any. The write will loop to the 1 addressable
// location, if the last addressable location is exceeded.
func (s *eeprom25XX256) Write(b []byte) (n int, err error) {
	n = 0
	for nextPageDistance := int(pageSize - (s.offset % pageSize)); n < len(b); nextPageDistance = 64 {
		if n+nextPageDistance > len(b) {
			nextPageDistance = len(b) - n
		}
		c, err := s.writeBlock(b[n : n+nextPageDistance])
		n += c
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (s *eeprom25XX256) writeBlock(b []byte) (n int, err error) {
	buf := make([]uint8, len(b)+3)
	buf[0] = instwrite
	binary.BigEndian.PutUint16(buf[1:], s.offset)
	copy(buf[3:], b)

	n = len(b)
	if _, err = s.bus.TransferAndReceiveByte(instwren); err != nil {
		return 0, err
	}

	if err = s.bus.TransferAndRecieveData(buf[:]); err != nil {
		return 0, err
	}

	if _, err = s.bus.TransferAndReceiveByte(instwrdi); err != nil {
		return 0, err
	}
	time.Sleep(twc)

	s.offset += uint16(n % maxAddress)
	return n, nil
}

// WriteAt writes len(b) bytes to the eeprom starting at byt offset off. It
// returns the number of bytes written and an error, if any. The read will loop to
// the 1 addressable location, if the last addressable location is exceeded.
func (s *eeprom25XX256) WriteAt(b []byte, off int64) (n int, err error) {
	if _, err = s.Seek(off, os.SEEK_SET); err != nil {
		return 0, err
	}

	return s.Write(b)
}
