package services

import (
    "backend/internal/database"
    "backend/internal/models"
    "context"
    "encoding/json"
    "errors"
)

var (
    ErrFormTemplateIdDoesNotExists = errors.New("form template does not exists for this user")
    ErrFormTemplateIdExists = errors.New("form template id already exists for this user")
)


func CreateFormTemplate(ctx context.Context, req *models.FormTemplate) error {
    db := database.GetDB()
    userID := ctx.Value("userID")

    // Check if form template already exists for this user
    var count int
    err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM form_templates WHERE form_template_id = $1 AND user_id = $2", req.FormTemplateID, userID)
    if err != nil {
        return err
    }
    if count > 0 {
        return ErrFormTemplateIdExists
    }

    // Insert new form template
	// Convert map to JSON for attributes
    fieldsJSON, err := json.Marshal(req.Fields)
    if err != nil {
        return err
    }


    query := `INSERT INTO form_templates (
        form_template_id, 
        user_id, 
        fields
    ) VALUES ($1, $2, $3)`
    _, err = db.ExecContext(ctx, query, 
        req.FormTemplateID, 
        userID, 
        fieldsJSON)

    return err
}


func GetFormTemplateById(ctx context.Context, templateID string) (*models.FormTemplate, error) {
    db := database.GetDB()
    userID := ctx.Value("userID")

    var template models.FormTemplate
	var fieldsJSON []byte

    query := `SELECT id, form_template_id, user_id, fields, created_at, updated_at 
              FROM form_templates WHERE form_template_id = $1 AND user_id = $2`
    
	err := db.QueryRowContext(ctx, query, templateID, userID).Scan(
			&template.ID,
			&template.FormTemplateID,
			&template.UserID,
			&fieldsJSON,
			&template.CreatedAt,
			&template.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Unmarshal attributes JSON
	if err := json.Unmarshal(fieldsJSON, &template.Fields); err != nil {
        return nil, err
    }

	return &template, nil
}


func GetFormTemplatesByUserId(ctx context.Context) ([]*models.FormTemplate, error) {
    db := database.GetDB()
    userID := ctx.Value("userID")

    var formTemplates []*models.FormTemplate
    query := `SELECT id, form_template_id, user_id, fields, created_at, updated_at
              FROM form_templates WHERE user_id = $1`

    rows, err := db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var formTemplate models.FormTemplate
        var fieldsJSON []byte

        err := rows.Scan(
            &formTemplate.ID,
            &formTemplate.FormTemplateID,
            &formTemplate.UserID,
            &fieldsJSON,
            &formTemplate.CreatedAt,
            &formTemplate.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }

        // Unmarshal the fields JSON into the map
        if err := json.Unmarshal(fieldsJSON, &formTemplate.Fields); err != nil {
            return nil, err
        }

        formTemplates = append(formTemplates, &formTemplate)
    }

    return formTemplates, nil
}


func DeleteFormTemplate(ctx context.Context, formTemplateID string) error {
    db := database.GetDB()
    userID := ctx.Value("userID")

    // Check if form template exists for this user
    var count int
    err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM form_templates WHERE form_template_id = $1 AND user_id = $2", formTemplateID, userID)
    if err != nil {
        return err
    }
    if count == 0 {
        return ErrFormTemplateIdDoesNotExists
    }

    // Delete the form template
    query := `DELETE FROM form_templates WHERE form_template_id = $1 AND user_id = $2`
    _, err = db.ExecContext(ctx, query, formTemplateID, userID)
    return err
}
