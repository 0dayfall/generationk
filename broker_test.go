package generationk

import (
	"testing"
	"time"
)

func TestNewBroker(t *testing.T) {

	broker := Broker{
		portfolio: &Portfolio{},
		callback:  nil,
		comission: nil,
	}

	if broker.callback != nil {
		t.Errorf("Empty broker could not be created")
	}

}

func TestBuyOrder(t *testing.T) {

	portfolio := NewPortfolio(1000)

	broker := Broker{
		portfolio: portfolio,
		callback:  nil,
		comission: nil,
	}

	cash, err := broker.SendOrder(Order{
		Asset:     "Any",
		direction: BuyOrder,
		orderType: 0,
		Time:      time.Time{},
		Price:     100,
		Qty:       10,
	}, nil)

	if err != nil {
		t.Error(err)
	}

	if cash != -1000 {
		t.Error("It should cost 1000 to buy the Any stock")
	}

}

func TestSellOrder(t *testing.T) {

	portfolio := NewPortfolio(1000)

	broker := Broker{
		portfolio: portfolio,
		callback:  nil,
		comission: nil,
	}

	cash, err := broker.SendOrder(Order{
		Asset:     "Any",
		direction: BuyOrder,
		orderType: 0,
		Time:      time.Time{},
		Price:     100,
		Qty:       10,
	}, nil)

	if err != nil {
		t.Error(err)
	}

	if portfolio.GetBalance()+cash == 0 {
		t.Error("We should not have any cash left")
	}

	cash, err = broker.SendOrder(Order{
		Asset:     "Any",
		direction: SellOrder,
		orderType: 0,
		Time:      time.Time{},
		Price:     200,
		Qty:       10,
	}, nil)

	if err != nil {
		t.Error(err)
	}

	if cash != 2000 {
		t.Error("We sold with a higher price and doubled our money")
	}

}
