package main

import (
	"backend-layout/internal/config"
	"backend-layout/internal/domain"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	fmt.Println("start seeder...")

	cfg, err := config.NewConfig(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	pgxConfig, err := pgxpool.ParseConfig(cfg.DB.DSN)

	if err != nil {
		fmt.Println("failed to parse pgx config: %w", err)
	}

	ctx := context.Background()
	dbpool, err := pgxpool.NewWithConfig(ctx, pgxConfig)

	if err != nil {
		fmt.Println("failed to create pgx pool: %w", err.Error())
	}

	seeder := new(dbpool)

	err = seeder.Run(ctx)

	if err != nil {
		fmt.Println("failed to seed %w", err.Error())
	}
	fmt.Println("completed..")
}

type seeder struct {
	conn *pgxpool.Pool
}

func new(conn *pgxpool.Pool) *seeder {
	return &seeder{
		conn: conn,
	}
}

func (s *seeder) Run(ctx context.Context) error {
	tx, err := s.conn.Begin(ctx)

	if err != nil {
		return err
	}

	err = s.seedBook(ctx, tx)

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	err = tx.Commit(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (s *seeder) seedBook(ctx context.Context, tx pgx.Tx) error {

	books := []domain.Book{
		{
			Title: "Statistika untuk Penelitian",
			Slug:  "statistika-untuk-penelitian-4",
			Author: domain.Author{
				Id: 1,
			},
			Publisher: domain.Publisher{
				Id: 1,
			},
			PublishYear: "2023",
			TotalPage:   220,
			Description: "lorem ipsum",
			Sku:         "208381357",
			Isbn:        "9786230142017",
			Price:       95000,
		},
		{
			Title: "Statistika Deskriptif: Teori, Aplikasi, dan Soal Pembahasan",
			Slug:  "statistika-deskriptif-teori-aplikasi-dan-soal-pembahasan",
			Author: domain.Author{
				Id: 1,
			},
			Publisher: domain.Publisher{
				Id: 1,
			},
			PublishYear: "2023",
			TotalPage:   300,
			Description: "lorem ipsum",
			Sku:         "208361329",
			Isbn:        "9786232059764",
			Price:       105000,
		},
		{
			Title: "Analisis Data Penelitian Kesehatan Memahami & Menggunakan SPSS",
			Slug:  "analisis-data-penelitian-kesehatan-memahami-menguna-spss",
			Author: domain.Author{
				Id: 1,
			},
			Publisher: domain.Publisher{
				Id: 1,
			},
			PublishYear: "2024",
			TotalPage:   350,
			Description: "lorem ipsum",
			Sku:         "208343747",
			Isbn:        "9786236913376",
			Price:       79000,
		},
		{
			Title: "Cepat Kuasai SPSS",
			Slug:  "cepat-kuasai-spss",
			Author: domain.Author{
				Id: 1,
			},
			Publisher: domain.Publisher{
				Id: 1,
			},
			PublishYear: "2023",
			TotalPage:   200,
			Description: "lorem ipsum",
			Sku:         "208340018",
			Isbn:        "9786231644367",
			Price:       55500,
		},
		{
			Title: "Panduan Lengkap SPSS 27",
			Slug:  "panduan-lengkap-spss-27",
			Author: domain.Author{
				Id: 1,
			},
			Publisher: domain.Publisher{
				Id: 1,
			},
			PublishYear: "2024",
			TotalPage:   300,
			Description: "lorem ipsum",
			Sku:         "723050784",
			Isbn:        "9786230055713",
			Price:       175000,
		},
		{
			Title: "Pengolahan dan Analisa Data Statistika dengan SPSS",
			Slug:  "pengolahan-dan-analisa-data-statistika-dengan-spss",
			Author: domain.Author{
				Id: 1,
			},
			Publisher: domain.Publisher{
				Id: 1,
			},
			PublishYear: "2024",
			TotalPage:   325,
			Description: "lorem ipsum",
			Sku:         "208335971",
			Isbn:        "9786230136528",
			Price:       136000,
		},
		{
			Title: "The Guide Book Of SPSS",
			Slug:  "the-guide-book-of-spss-cara-mudah-dan-cepat-mengolah-data",
			Author: domain.Author{
				Id: 1,
			},
			Publisher: domain.Publisher{
				Id: 1,
			},
			PublishYear: "2023",
			TotalPage:   400,
			Description: "lorem ipsum",
			Sku:         "208310800",
			Isbn:        "9786234008937",
			Price:       53550,
		},
		{
			Title: "Statistik Deskriptif",
			Slug:  "statistik-deskriptif-1",
			Author: domain.Author{
				Id: 1,
			},
			Publisher: domain.Publisher{
				Id: 1,
			},
			PublishYear: "2015",
			TotalPage:   240,
			Description: "lorem ipsum",
			Sku:         "208312164",
			Isbn:        "9786235690407",
			Price:       35100,
		},
		{
			Title: "Statistik Untuk Penelitian Psikologi Dengan SPSS",
			Slug:  "statistik-untuk-penelitian-psikologi-dengan-spss",
			Author: domain.Author{
				Id: 1,
			},
			Publisher: domain.Publisher{
				Id: 1,
			},
			PublishYear: "2023",
			TotalPage:   250,
			Description: "lorem ipsum",
			Sku:         "208298701",
			Isbn:        "9786233729369",
			Price:       137700,
		},
		{
			Title: "Statistik Multivariat",
			Slug:  "statistik-multivariat",
			Author: domain.Author{
				Id: 1,
			},
			Publisher: domain.Publisher{
				Id: 1,
			},
			PublishYear: "2023",
			TotalPage:   320,
			Description: "lorem ipsum",
			Sku:         "208123719",
			Isbn:        "9786024254780",
			Price:       219600,
		},
		{
			Title: "The Master Book Of SPSS",
			Slug:  "the-master-book-of-spss",
			Author: domain.Author{
				Id: 1,
			},
			Publisher: domain.Publisher{
				Id: 1,
			},
			PublishYear: "2025",
			TotalPage:   450,
			Description: "lorem ipsum",
			Sku:         "208033038",
			Isbn:        "9786237324362",
			Price:       250000,
		},
		{
			Title: "Pengolahan Data Kesehatan Dengan SPSS",
			Slug:  "pengolahan-data-kesehatan-dengan-spss",
			Author: domain.Author{
				Id: 1,
			},
			Publisher: domain.Publisher{
				Id: 1,
			},
			PublishYear: "2023",
			TotalPage:   200,
			Description: "lorem ipsum",
			Sku:         "208018288",
			Isbn:        "9786025375927",
			Price:       37800,
		},
		{
			Title: "Mahir Statistik Multivariat dengan SPSS",
			Slug:  "mahir-statistik-multivariat-dengan-spss",
			Author: domain.Author{
				Id: 1,
			},
			Publisher: domain.Publisher{
				Id: 1,
			},
			PublishYear: "2024",
			TotalPage:   300,
			Description: "lorem ipsum",
			Sku:         "718051071",
			Isbn:        "9786020477169",
			Price:       86750,
		},
	}

	batch := &pgx.Batch{}

	for _, v := range books {
		batch.Queue(`INSERT INTO books (title,
    slug,
    author_id,
    publisher_id,
    publish_year,
    total_page,
    description,
    sku,
    isbn,
    price) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, v.Title, v.Slug, v.Author.Id, v.Publisher.Id, v.PublishYear, v.TotalPage, v.Description, v.Sku, v.Isbn, v.Price)
	}

	br := tx.SendBatch(ctx, batch)

	defer br.Close()

	return br.Close()
}
