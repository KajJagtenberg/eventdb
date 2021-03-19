package store

type StoreService struct {
	storage *Storage
}

func (s *StoreService) Add(req *AddRequest, stream EventStore_AddServer) error {
	events, err := s.storage.Add(req)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := stream.Send(event); err != nil {
			return err
		}
	}

	return nil
}

func (s *StoreService) Get(req *GetRequest, stream EventStore_GetServer) error {
	events, err := s.storage.Get(req)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := stream.Send(event); err != nil {
			return err
		}
	}

	return nil
}

func (s *StoreService) Log(req *LogRequest, stream EventStore_LogServer) error {
	events, err := s.storage.Log(req)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := stream.Send(event); err != nil {
			return err
		}
	}

	return nil
}

func NewStoreService(storage *Storage) *StoreService {
	return &StoreService{storage}
}
