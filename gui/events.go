package main

import (
	"os"

	"github.com/therecipe/qt/core"
)

func init() {
	CustomTableModel_QmlRegisterType2("CustomQmlTypes", 1, 0, "CustomTableModel")
	BridgeTemplate_QmlRegisterType2("CustomQmlTypes", 1, 0, "BridgeTemplate")
}

const (
	Time = int(core.Qt__UserRole) + 1<<iota
	Event
	Value
)

type EventTableItem struct {
	time  string
	event string
	value string
}

type CustomTableModel struct {
	core.QAbstractTableModel

	_ func() `constructor:"init"`

	_ func()                                        `signal:"remove,auto"`
	_ func(item []*core.QVariant)                   `signal:"add,auto"`
	_ func(time string, event string, value string) `signal:"edit,auto"`

	modelData []EventTableItem
}

type BridgeTemplate struct {
	core.QObject

	_ func() `signal:"quitBtnClicked,auto"`
}

func (m *CustomTableModel) init() {
	m.modelData = []EventTableItem{{"0", "RPN LSB", "0"}}

	m.ConnectRoleNames(m.roleNames)
	m.ConnectRowCount(m.rowCount)
	m.ConnectColumnCount(m.columnCount)
	m.ConnectData(m.data)
}

func (m *CustomTableModel) roleNames() map[int]*core.QByteArray {
	return map[int]*core.QByteArray{
		Time:  core.NewQByteArray2("Time", -1),
		Event: core.NewQByteArray2("Event", -1),
		Value: core.NewQByteArray2("Value", -1),
	}
}

func (m *CustomTableModel) rowCount(*core.QModelIndex) int {
	return len(m.modelData)
}

func (m *CustomTableModel) columnCount(*core.QModelIndex) int {
	return 3
}

func (m *CustomTableModel) data(index *core.QModelIndex, role int) *core.QVariant {
	item := m.modelData[index.Row()]
	switch role {
	case Time:
		return core.NewQVariant1(item.time)
	case Event:
		return core.NewQVariant1(item.event)
	case Value:
		return core.NewQVariant1(item.value)
	}
	return core.NewQVariant()
}

func (m *CustomTableModel) remove() {
	if len(m.modelData) == 0 {
		return
	}
	m.BeginRemoveRows(core.NewQModelIndex(), len(m.modelData)-1, len(m.modelData)-1)
	m.modelData = m.modelData[:len(m.modelData)-1]
	m.EndRemoveRows()
}

func (m *CustomTableModel) add(item []*core.QVariant) {
	m.BeginInsertRows(core.NewQModelIndex(), len(m.modelData), len(m.modelData))
	m.modelData = append(m.modelData, EventTableItem{item[0].ToString(), item[1].ToString(), item[2].ToString()})
	m.EndInsertRows()
}

func (m *CustomTableModel) edit(time string, event string, value string) {
	if len(m.modelData) == 0 {
		return
	}
	m.modelData[len(m.modelData)-1] = EventTableItem{time, event, value}
	m.DataChanged(m.Index(len(m.modelData)-1, 0, core.NewQModelIndex()), m.Index(len(m.modelData)-1, 1, core.NewQModelIndex()), []int{Time, Event, Value})
}

func (b *BridgeTemplate) quitBtnClicked() {
	os.Exit(0)
}
