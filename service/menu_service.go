package service

import (
	"context"
	"log"

	pb "github.com/Mubinabd/reservation_service/genproto"
	"github.com/Mubinabd/reservation_service/storage"
	"github.com/google/uuid"
)

type MenuService struct {
	stg storage.StorageI
	pb.UnimplementedMenuServiceServer
}

// mustEmbedUnimplementedMenuServiceServer implements genproto.MenuServiceServer.
// func (menu *MenuService) mustEmbedUnimplementedMenuServiceServer() {
// 	panic("unimplemented")
// }

func NewMenuService(stg storage.StorageI) *MenuService {
	return &MenuService{stg: stg}
}

func (menu *MenuService) CreateMenu(ctx context.Context, req *pb.Menu) (*pb.Void, error) {
	id := uuid.NewString()
	req.Id = id
	_, err := menu.stg.Menu().Create(req)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ps *MenuService) GetMenu(ctx context.Context, menu *pb.ById) (*pb.Menu, error) {
	res, err := ps.stg.Menu().GetById((*pb.ById)(menu))
	if err != nil {
		log.Fatal(err.Error())
	}

	return res, err
}

func (ps *MenuService) GetMenus(ctx context.Context, menu *pb.MenuFilter) (*pb.Menus, error) {
	res, err := ps.stg.Menu().GetAll((*pb.MenuFilter)(menu))
	if err != nil {
		return nil, err
	}
	return res, err
}

func (ps *MenuService) UpdateMenu(ctx context.Context, menu *pb.Menu) (*pb.Void, error) {
	_, err := ps.stg.Menu().Update(menu)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ps *MenuService) DeleteMenu(ctx context.Context, menu *pb.ById) (*pb.Void, error) {
	_, err := ps.stg.Menu().Delete(menu)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
