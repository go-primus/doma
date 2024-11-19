package asynctask

type Store struct {
	state any
}

func (s *Store) Dispatch(state any) {
	s.state = state

	// notify
}

type ReducerStore struct {
	//

}

func (s *ReducerStore) Signal() {

}
