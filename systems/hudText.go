package systems

import (
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"image/color"
)

// Text entity for text printed to screen
type Text struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
}

// HUDTextMessage updates hud text from messages
type HUDTextMessage struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.MouseComponent
	Line1, Line2, Line3, Line4 string
}

const HUDTextMessageType string = "HUDTextMessage"

func (HUDTextMessage) Type() string {
	return HUDTextMessageType
}

type HUDMoneyMessage struct {
	CurrentAmount    int
	IncomePerTenHour int
}

const HUDMoneyMessageType string = "HUDMoneyMessage"

func (HUDMoneyMessage) Type() string {
	return HUDMoneyMessageType
}

type HUDTextEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*common.MouseComponent
	Line1, Line2, Line3, Line4 string
}

// HUDTextSystem system prints text to HUD
type HUDTextSystem struct {
	text1, text2  Text
	text3, text4  Text
	currentTime   Text
	incomePerText Text
	moneyText     Text
	gameTime      float32
	updateMoney   bool //Let us know if amount of money has been updated
	moneyAmount   int  //keeps track of amount of money to display
	incomePer     int
	mouse         common.MouseComponent

	entities []HUDTextEntity
}

func (hud *HUDTextSystem) New(w *ecs.World) {

	setDefaultTextValues(w, &hud.text1, "Nothing Selected!", engo.WindowHeight()-200)
	setDefaultTextValues(w, &hud.text2, "hover on a town", engo.WindowHeight()-180)
	setDefaultTextValues(w, &hud.text3, "to get info", engo.WindowHeight()-160)
	setDefaultTextValues(w, &hud.text4, "about it", engo.WindowHeight()-140)
	setDefaultTextValues(w, &hud.currentTime, "Time(hours): 0", engo.WindowHeight()-100)
	setDefaultTextValues(w, &hud.incomePerText, "Income every 10hr: $0", engo.WindowHeight()-60)
	setDefaultTextValues(w, &hud.moneyText, "Available funds: $0", engo.WindowHeight()-40)

	engo.Mailbox.Listen(HUDTextMessageType, func(m engo.Message) {
		msg, ok := m.(HUDTextMessage)
		if !ok {
			return
		}
		for _, system := range w.Systems() {
			switch sys := system.(type) {
			case *common.MouseSystem:
				sys.Add(&msg.BasicEntity, &msg.MouseComponent,
					&msg.SpaceComponent, nil)
			case *HUDTextSystem:
				sys.Add(&msg.BasicEntity, &msg.SpaceComponent,
					&msg.MouseComponent, msg.Line1, msg.Line2,
					msg.Line3, msg.Line4)
			}
		}
	})

	engo.Mailbox.Listen(HUDMoneyMessageType, func(m engo.Message) {
		msg, ok := m.(HUDMoneyMessage)
		if !ok {
			return
		}
		hud.moneyAmount = msg.CurrentAmount
		hud.incomePer = msg.IncomePerTenHour
		hud.updateMoney = true
	})
}

func setDefaultTextValues(w *ecs.World, text *Text,
	textContent string, yPos float32) {
	fnt := &common.Font{
		URL:  "go.ttf",
		FG:   color.Black,
		Size: 24,
	}
	err := fnt.CreatePreloaded()
	if err != nil {
	}
	text.BasicEntity = ecs.NewBasic()
	text.RenderComponent.Drawable = common.Text{
		Font: fnt,
		Text: textContent,
	}

	text.SetShader(common.TextHUDShader)
	text.RenderComponent.SetZIndex(103)
	text.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 0, Y: yPos},
		Width:    200,
		Height:   200,
	}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&text.BasicEntity, &text.RenderComponent, &text.SpaceComponent)
		}
	}
}

func (hud *HUDTextSystem) Add(b *ecs.BasicEntity, s *common.SpaceComponent, m *common.MouseComponent, l1, l2, l3, l4 string) {
	hud.entities = append(hud.entities, HUDTextEntity{b, s, m, l1, l2, l3, l4})
}

func (hud *HUDTextSystem) Update(dt float32) {
	hud.gameTime += dt
	isHovering := false
	for _, e := range hud.entities {
		if e.MouseComponent.Hovered {
			txt := hud.text1.RenderComponent.Drawable.(common.Text)
			txt.Text = e.Line1
			hud.text1.RenderComponent.Drawable = txt
			txt = hud.text2.RenderComponent.Drawable.(common.Text)
			txt.Text = e.Line2
			hud.text2.RenderComponent.Drawable = txt
			txt = hud.text3.RenderComponent.Drawable.(common.Text)
			txt.Text = e.Line3
			hud.text3.RenderComponent.Drawable = txt
			txt = hud.text4.RenderComponent.Drawable.(common.Text)
			txt.Text = e.Line4
			hud.text4.RenderComponent.Drawable = txt
			isHovering = true

		}
	}
	if !isHovering {
		txt := hud.text1.RenderComponent.Drawable.(common.Text)
		txt.Text = "Hover on town for info"
		hud.text1.RenderComponent.Drawable = txt
		txt = hud.text2.RenderComponent.Drawable.(common.Text)
		txt.Text = ""
		hud.text2.RenderComponent.Drawable = txt
		txt = hud.text3.RenderComponent.Drawable.(common.Text)
		txt.Text = "Press F1,F2 to add"
		hud.text3.RenderComponent.Drawable = txt
		txt = hud.text4.RenderComponent.Drawable.(common.Text)
		txt.Text = "or remove highways"
		hud.text4.RenderComponent.Drawable = txt
	}
	txt := hud.currentTime.RenderComponent.Drawable.(common.Text)
	txt.Text = fmt.Sprintf("Time(hours): %.0f", hud.gameTime)
	hud.currentTime.RenderComponent.Drawable = txt

	txt = hud.incomePerText.RenderComponent.Drawable.(common.Text)
	txt.Text = fmt.Sprintf("Income every 10hr: $%v", hud.incomePer)
	hud.incomePerText.RenderComponent.Drawable = txt

	if hud.updateMoney {
		txt := hud.moneyText.RenderComponent.Drawable.(common.Text)
		txt.Text = fmt.Sprintf("Funds Available: $%v", hud.moneyAmount)
		hud.moneyText.RenderComponent.Drawable = txt
	}
}

func (hud *HUDTextSystem) Remove(entity ecs.BasicEntity) {}
