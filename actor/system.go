package actor

import (
  "secure-route-api/city"
  "secure-route-api/crime"
  "time"

  "github.com/AsynkronIT/protoactor-go/actor"
)

type System struct {
  PID      *actor.PID
  Children map[int]*actor.PID
  cities   map[int]city.City
  crimes   *crime.Crimes
}

func NewSystem(cities map[int]city.City, crimes *crime.Crimes) System {
  system := &System{cities: cities, crimes: crimes}
  context := actor.EmptyRootContext
  props := actor.PropsFromProducer(func() actor.Actor {
    return system
  })
  system.PID = context.Spawn(props)
  return *system
}

func (system *System) TellSync(msg interface{}) (interface{}, error) {
  context := actor.EmptyRootContext
  future := context.RequestFuture(system.PID, msg, 7*time.Second)
  return future.Result()
}

func (system *System) Tell(msg interface{}) {
  context := actor.EmptyRootContext
  context.Send(system.PID, msg)
}

func (system *System) Receive(context actor.Context) {
  switch msg := context.Message().(type) {
  case *actor.Started:
    system.Children = make(map[int]*actor.PID)
    for _, c := range system.cities {
      props := actor.PropsFromProducer(func() actor.Actor {
        return &City{
          //TODO this crimes should be by city.
          Crimes: system.crimes,
        }
      })
      system.Children[c.ID] = context.Spawn(props)
    }

  case Route:
    child := system.Children[msg.CityID]
    future := actor.EmptyRootContext.RequestFuture(child, msg, 70*time.Second)
    r, _ := future.Result()
    context.Respond(r)
  }
}
