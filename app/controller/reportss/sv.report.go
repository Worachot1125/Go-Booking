package reportss

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"errors"
	"fmt"
)

func (s *Service) Create(ctx context.Context, req request.CreateReport) (*response.ReportResponse, error) {
	m := &model.Report{
		UserID:      req.UserID,
		RoomID:      req.RoomID,
		Description: req.Description,
	}

	_, err := s.db.NewInsert().Model(m).Returning("*").Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	resp := &response.ReportResponse{
		ID:          m.ID,
		User_ID:     m.UserID,
		Room_ID:     m.RoomID,
		Name_user:   m.Name_user,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}

	return resp, nil
}

func (s *Service) Update(ctx context.Context, req request.UpdateReport, id request.GetByIDReport) (*response.ReportResponse, error) {
	m := &model.Report{
		ID:          id.ID,
		UserID:      req.UserID,
		RoomID:      req.RoomID,
		Description: req.Description,
	}
	m.SetUpdateNow()

	res, err := s.db.NewUpdate().
		Model(m).
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update report: %w", err)
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return nil, errors.New("report not found")
	}

	resp := &response.ReportResponse{
		ID:          m.ID,
		User_ID:     m.UserID,
		Room_ID:     m.RoomID,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}

	return resp, nil
}

func (s *Service) List(ctx context.Context) ([]response.ReportResponse, error) {
    var reports []model.Report

    err := s.db.NewSelect().
        Model(&reports).
        Column("rep.*").
        ColumnExpr(`COALESCE("user".first_name || ' ' || "user".last_name, '') AS name_user`).
        Join(`JOIN users AS "user" ON "user".id = rep.user_id::uuid AND "user".deleted_at IS NULL`).
        Where(`rep.deleted_at IS NULL`).
        Scan(ctx)
    if err != nil {
        return nil, err
    }

    // map ออก response โดยใช้ r.Name_user ที่สแกนมาจาก SQL
    res := make([]response.ReportResponse, 0, len(reports))
    for _, r := range reports {
        res = append(res, response.ReportResponse{
            ID:          r.ID,
            User_ID:     r.UserID,
            Name_user:   r.Name_user,
            Room_ID:     r.RoomID,
            Description: r.Description,
            CreatedAt:   r.CreatedAt,
            UpdatedAt:   r.UpdatedAt,
        })
    }
    return res, nil
}


func (s *Service) Get(ctx context.Context, id string) (*response.ReportResponse, error) {
    var report model.Report

    err := s.db.NewSelect().
        Model(&report).
        Column("rep.*").
        ColumnExpr(`COALESCE(u.first_name || ' ' || u.last_name, '') AS name_user`).
        Join(`JOIN users AS u ON u.id = rep.user_id::uuid AND u.deleted_at IS NULL`).
        Where("rep.deleted_at IS NULL").
        Where("rep.id = ?", id).
        Scan(ctx)
    if err != nil {
        return nil, err
    }

    res := &response.ReportResponse{
        ID:          report.ID,
        User_ID:     report.UserID,
        Name_user:   report.Name_user,
        Room_ID:     report.RoomID,
        Description: report.Description,
        CreatedAt:   report.CreatedAt,
        UpdatedAt:   report.UpdatedAt,
    }
    return res, nil
}

func (s *Service) Delete(ctx context.Context, id request.GetByIDReport) error {
	ex, err := s.db.NewSelect().
		Table("reports").
		Where("id = ?", id.ID).
		Where("deleted_at IS NULL").
		Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check report existence: %w", err)
	}
	if !ex {
		return errors.New("report not found")
	}

	_, err = s.db.NewDelete().
		Model((*model.Report)(nil)).
		Where("id = ?", id.ID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete report: %w", err)
	}
	return nil
}
