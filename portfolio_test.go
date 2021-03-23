package generationk

import (
	"errors"
	"testing"
)

func TestNewPortfolio(t *testing.T) {

	portfolio := NewPortfolio()

	if portfolio.GetBalance() != 0 {
		t.Errorf("Error creating portfolio with zero balance")
	}

}

func TestIsOwning(t *testing.T) {

	portfolio := NewPortfolio()

	portfolio.AddHolding(Holding{assetName: "holding"})

	if portfolio.IsOwning("holding") == false {
		t.Errorf("Could not find holding")
	}

	if portfolio.IsOwning("holding2") == true {
		t.Errorf("holding2 is not owned")
	}

	portfolio.RemoveHolding(Holding{assetName: "holding"})

	if portfolio.IsOwning("holding") == true {
		t.Errorf("Holding was removed")
	}

}

func TestAddRemoveHolding(t *testing.T) {

	portfolio := NewPortfolio()

	portfolio.AddHolding(Holding{assetName: "holding"})

	if portfolio.IsOwning("holding") == false {
		t.Errorf("Could not find holding")
	}

	portfolio.AddHolding(Holding{assetName: "holding2"})

	if len(portfolio.holdings) != 2 {
		t.Errorf("Holding has wrong length")
	}

	portfolio.RemoveHolding(Holding{assetName: "holding"})

	if portfolio.IsOwning("holding") == true {
		t.Errorf("Holding was removed")
	}

}

func TestBalance(t *testing.T) {

	portfolio := NewPortfolio()
	portfolio.SetBalance(10000)

	err := portfolio.subtractFromBalance(300)
	if err != nil {
		t.Error(err)
	}

	if portfolio.GetBalance() != 9700 {
		t.Errorf("Wrong balance")
	}

	err = portfolio.subtractFromBalance(9800)

	if !errors.Is(&NegativeBalanceErr{}, err) {
		t.Error("Did not get negative balance error")
	}

	err = portfolio.subtractFromBalance(9700)
	if err != nil {
		t.Error("Something went wrong")
	}

	portfolio.addToBalance(300)

	if portfolio.GetBalance() != 300 {
		t.Errorf("Wrong balance")
	}

}
