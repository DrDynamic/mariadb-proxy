package server

import (
	"fmt"
	"mschon/dbproxy/readers"
	"mschon/dbproxy/tcp/mariadb/packets"
)

const PacketTypeErrPacket packets.PacketType = "ErrPacket"

type ErrPacket struct {
	packets.Packet
	ErrorCode uint16

	ReportType     ReportType
	ProgressReport ProgressReport
	SqlError       SqlError
	GeneralError   string
}

type ReportType int

const (
	REPORT_TYPE_PROGRESS = iota
	REPORT_TYPE_SQL
	REPORT_TYPE_GENERAL
)

type ProgressReport struct {
	Stage        uint8
	MaxStage     uint8
	Progress     uint32
	ProgressInfo string
}

type SqlError struct {
	SqlState string
	Message  string
}

func NewErrPacket(packet *packets.Packet) *ErrPacket {
	reader := readers.NewLeByteReader(packet.Body)
	code := reader.PopUInt16()

	var reportType ReportType
	var progressReport ProgressReport
	var sqlError SqlError
	var generalError string
	if code == 0xFFFF {
		reportType = REPORT_TYPE_PROGRESS
		stage := reader.PopUInt8()
		maxStage := reader.PopUInt8()
		progress := reader.PopUInt24()
		info := reader.PopLengthEncodedString()
		progressReport = ProgressReport{
			Stage:        stage,
			MaxStage:     maxStage,
			Progress:     progress,
			ProgressInfo: info,
		}
	} else if packet.Body[3] == '#' {
		reportType = REPORT_TYPE_SQL
		reader.PopBytes(1) // state marker
		sqlState := reader.PopString(5)
		message := reader.PopEofEncodedString()
		sqlError = SqlError{
			SqlState: sqlState,
			Message:  message,
		}
	} else {
		reportType = REPORT_TYPE_GENERAL
		generalError = reader.PopEofEncodedString()
	}

	packet.Type = PacketTypeErrPacket

	return &ErrPacket{
		Packet:         *packet,
		ErrorCode:      code,
		ReportType:     reportType,
		ProgressReport: progressReport,
		SqlError:       sqlError,
		GeneralError:   generalError,
	}
}

func (packet ErrPacket) GetType() packets.PacketType {
	return packet.Type
}

func (packet ErrPacket) GetSequence() byte {
	return packet.Sequence
}

func (packet *ErrPacket) String() string {
	switch packet.ReportType {
	case REPORT_TYPE_PROGRESS:
		return fmt.Sprintf("[ErrPacket PROGRESS code=%d max=%d stage=%d progress=%d info=%s]",
			packet.ErrorCode,
			packet.ProgressReport.MaxStage,
			packet.ProgressReport.Stage,
			packet.ProgressReport.Progress,
			packet.ProgressReport.ProgressInfo)
	case REPORT_TYPE_SQL:
		return fmt.Sprintf("[ErrPacket SQL code=%d state=%s message=%s]",
			packet.ErrorCode,
			packet.SqlError.SqlState,
			packet.SqlError.Message)
	case REPORT_TYPE_GENERAL:
		return fmt.Sprintf("[ErrPacket SQL code=%d message=%s]",
			packet.ErrorCode,
			packet.GeneralError)
	}
	return packet.Packet.String()
}
