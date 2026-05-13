package game

import "errors"

var (
    ErrInvalidCell    = errors.New("invalid cell index")
    ErrNotYourCell    = errors.New("this is not your cell")
    ErrHouseCell      = errors.New("cannot plant in a house")
    ErrCellOccupied   = errors.New("cell already has a fruit")
    ErrNoFruit        = errors.New("no fruit in this cell")
    ErrNotRipe        = errors.New("fruit is not ripe yet")
    ErrNotYourFruit   = errors.New("this is not your fruit")
    ErrCellTaken      = errors.New("cell is already taken by another team")
    ErrNotEnoughMoney = errors.New("not enough money")
    ErrNotEnoughFruit = errors.New("not enough fruit in inventory")
    ErrQuestNotFound  = errors.New("quest not found")
    ErrInvalidTeam    = errors.New("invalid team id")
)
