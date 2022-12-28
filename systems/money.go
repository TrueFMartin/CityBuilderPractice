package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
)

type MoneySystem struct {
	amount                int
	towns, cities, metros int
	officers              int
	elapsed               float32
}

type CityType int

const (
	CityTypeNew = iota
	CityTypeTown
	CityTypeCity
	CityTypeMetro
)

func (m *MoneySystem) New(w *ecs.World) {
	engo.Mailbox.Listen(CityUpdateMessageType, func(msg engo.Message) {
		update, ok := msg.(CityUpdateMessage)
		if !ok {
			return
		}
		//Removes old city from count
		oldRemove := func() {
			switch update.Old {
			case CityTypeTown:
				m.towns--
			case CityTypeCity:
				m.cities--
			case CityTypeMetro:
				m.metros--
			}
		}
		switch update.New {
		case CityTypeNew:
			m.towns++
		case CityTypeTown:
			m.towns++
			oldRemove()
		case CityTypeCity:
			m.cities--
			oldRemove()
		case CityTypeMetro:
			m.metros--
			oldRemove()
		}
	})

	engo.Mailbox.Listen(AddOfficerMessageType, func(engo.Message) {
		m.officers++
	})
}

func (m *MoneySystem) Update(dt float32) {
	m.elapsed += dt
	if m.elapsed > 10 {
		m.amount += m.towns*100 + m.cities*500 + m.metros*1000
		m.amount -= m.officers * 20
		engo.Mailbox.Dispatch(HUDMoneyMessage{
			Amount: m.amount,
		})
		m.elapsed = 0
	}
}

func (m *MoneySystem) Remove(e ecs.BasicEntity) {}

type CityUpdateMessage struct {
	Old, New CityType
}

const CityUpdateMessageType string = "CityUpdateMessage"

func (CityUpdateMessage) Type() string {
	return CityUpdateMessageType
}

type AddOfficerMessage struct{}

const AddOfficerMessageType string = "AddOfficerMessage"

func (AddOfficerMessage) Type() string {
	return AddOfficerMessageType
}
