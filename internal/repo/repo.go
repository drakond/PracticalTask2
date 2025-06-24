package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) *Repository {
	return &Repository{conn: conn}
}

// User CRUD
func (r *Repository) CreateUser(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, created_at`
	return r.conn.QueryRow(ctx, query, user.Username, user.Password).Scan(&user.ID, &user.CreatedAt)
}

func (r *Repository) GetUserByID(ctx context.Context, id int) (*User, error) {
	user := &User{}
	query := `SELECT id, username, password, created_at FROM users WHERE id = $1`
	err := r.conn.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	user := &User{}
	query := `SELECT id, username, password, created_at FROM users WHERE username = $1`
	err := r.conn.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (r *Repository) DeleteUser(ctx context.Context, id int) error {
	cmd, err := r.conn.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}

// Task CRUD
func (r *Repository) CreateTask(ctx context.Context, task *Task) error {
	query := `INSERT INTO tasks (user_id, title, description, status) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return r.conn.QueryRow(ctx, query, task.UserID, task.Title, task.Description, task.Status).Scan(&task.ID, &task.CreatedAt)
}

func (r *Repository) GetTaskByID(ctx context.Context, id int) (*Task, error) {
	task := &Task{}
	query := `SELECT id, user_id, title, description, status, created_at FROM tasks WHERE id = $1`
	err := r.conn.QueryRow(ctx, query, id).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("task not found")
	}
	return task, err
}

func (r *Repository) GetAllTasks(ctx context.Context) ([]*Task, error) {
	rows, err := r.conn.Query(ctx, `SELECT id, user_id, title, description, status, created_at FROM tasks`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*Task{}
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *Repository) GetTasksByUserID(ctx context.Context, userID int) ([]*Task, error) {
	rows, err := r.conn.Query(ctx, `SELECT id, user_id, title, description, status, created_at FROM tasks WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*Task{}
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *Repository) GetTasksByUsername(ctx context.Context, username string) ([]*Task, error) {
	query := `SELECT t.id, t.user_id, t.title, t.description, t.status, t.created_at FROM users u LEFT JOIN tasks t ON t.user_id = u.id WHERE u.username = $1`
	rows, err := r.conn.Query(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*Task{}
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *Repository) UpdateTask(ctx context.Context, id int, title, description, status string) error {
	cmd, err := r.conn.Exec(ctx,
		`UPDATE tasks SET title = $1, description = $2, status = $3 WHERE id = $4`,
		title, description, status, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("task not found")
	}
	return nil
}

func (r *Repository) DeleteTask(ctx context.Context, id int) error {
	cmd, err := r.conn.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("task not found")
	}
	return nil
}
