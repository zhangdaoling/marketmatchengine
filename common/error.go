package common

import "errors"

var ErrIrregularData = errors.New("irregular data")
var ErrUnexpectedEOF = errors.New("unexpected EOF")
var ErrTooLarge = errors.New("bytes.Buffer: too large")
var ErrExpectedEOF = errors.New("Expected EOF")

var ErrQueueType = errors.New("irregular queue type")
var ErrEngineCheckSum = errors.New("engine checksum err")
