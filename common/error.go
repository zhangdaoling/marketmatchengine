package common

import "errors"

var ErrIrregularData = errors.New("irregular data")
var ErrUnexpectedEOF = errors.New("unexpected EOF")
var ErrTooLarge = errors.New("bytes.Buffer: too large")
var ErrExpectedEOF = errors.New("Expected EOF")
var ErrQueueType = errors.New("irregular queue type")
var ErrEngineCheckSum = errors.New("engine checksum err")

var ErrSymbol = errors.New("symbol err")

var ErrQuotationAmount = errors.New("amount is error after sub")

var ErrNotExist = errors.New("price is not exist int quotation slice")

var ErrAmount = errors.New("amount < 0")
