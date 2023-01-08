package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
)

type MoneySystem struct {
	amount                int
	incomePer             int
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
			case pathTownLoc:
				m.towns--
			case pathCityLoc:
				m.cities--
			case pathMetroLoc:
				m.metros--
			}
		}
		switch update.New {
		case townLoc:
			m.towns++
		case pathTownLoc:
			m.towns++
			oldRemove()
		case pathCityLoc:
			m.cities--
			oldRemove()
		case pathMetroLoc:
			m.metros--
			oldRemove()
		}
	})
	//On police officer add, update hudText
	engo.Mailbox.Listen(AddOfficerMessageType, func(engo.Message) {
		m.officers++
	})
	//Update amount of money every time road is built/removed
	engo.Mailbox.Listen(RoadCostMessageType, func(msg engo.Message) {
		update, ok := msg.(RoadCostMessage)
		if !ok {
			return
		}
		m.amount += update.Amount
	})
}

func (m *MoneySystem) Update(dt float32) {
	m.elapsed += dt //used to increase funds every 10 seconds
	//Combine income from all units minus officers
	m.incomePer = m.towns*100 + m.cities*500 + m.metros*1000 - m.officers*20
	//every 10 seconds, apply income change
	if m.elapsed > 10 {
		m.amount += m.incomePer
		m.elapsed = 0
	}
	//update HUD money display
	engo.Mailbox.Dispatch(HUDMoneyMessage{
		CurrentAmount:    m.amount,
		IncomePerTenHour: m.incomePer,
	})
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
