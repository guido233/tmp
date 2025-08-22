package modbus

import (
	"go-app/conf"
	"go-app/logger"
	"time"

	modbus "github.com/simonvetter/modbus"
)

func RtuModbusClient(config conf.Config) (*modbus.ModbusClient, error) {
	logger.Infof("RTUModbus start")
	parity := 0
	switch config.Rtumodbus.Parity {
	case "N":
		parity = 0
	case "E":
		parity = 1
	case "O":
		parity = 2
	}
	// Modbus RTU/ASCII
	// for an RTU (serial) device/bus
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:      "rtu://" + config.Rtumodbus.Device,
		Speed:    uint(config.Rtumodbus.BaudRate), // default
		DataBits: uint(config.Rtumodbus.DataBits), // default, optional
		Parity:   uint(parity),                    // default, optional
		StopBits: uint(config.Rtumodbus.StopBits), // default if no parity, optional
		Timeout:  1 * time.Second,
	})
	if err != nil {
		logger.Fatalf("rtumodbus connect error!")
		return nil, err
	}
	err = client.SetEncoding(modbus.BIG_ENDIAN, modbus.LOW_WORD_FIRST)
	if err != nil {
		logger.Fatalf("rtumodbus SetEncoding error!")
		return nil, err
	}
	err = client.Open()
	if err != nil {
		logger.Fatalf("rtumodbus connect error!")
		return nil, err
	}

	logger.Infof("rtumodbus connect to " + config.Rtumodbus.Device + " successful")

	return client, nil
}
