package services

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sakithb/hcblk-server/internal/models"
)

type ListingService struct {
	DB *sqlx.DB
}

func (s *ListingService) PopulateListing(l *models.Listing) error {
	l.PhoneNos = strings.Split(l.PhoneNosRaw, ",")

	imgs, err := os.ReadDir("./assets/dist/listings/" + l.Id + "/")
	if err != nil {
		return err
	}

	for _, img := range imgs {
		l.Images = append(l.Images, img.Name())
	}

	return nil
}

func (s *ListingService) GetListingById(id string) (*models.Listing, error) {
	l := models.Listing{}
	err := s.DB.Get(&l, `
		SELECT 
		l.id, l.description, l.price, l.mileage, l.used, l.phone_nos, l.listed_at,
		bike.id "bike.id", bike.model "bike.model", bike.brand "bike.brand", bike.category "bike.category", bike.year "bike.year", bike.engine_capacity "bike.engine_capacity",
		city.id "city.id", city.city "city.city", city.district "city.district", city.province "city.province",
		seller.id "seller.id", seller.first_name "seller.first_name", seller.last_name "seller.last_name", seller.email "seller.email", seller.joined_at "seller.joined_at"
		FROM listings AS l
		INNER JOIN bikes as bike ON l.bike_id = bike.id
		INNER JOIN cities AS city ON l.city_id = city.id
		INNER JOIN users AS seller ON l.seller_id = seller.id
		WHERE l.id = ?
	`, id)

	if err != nil {
		return nil, err
	}

	err = s.PopulateListing(&l)

	return &l, err
}

func (s *ListingService) GetListingsBySeller(sellerId string) ([]*models.Listing, error) {
	ls := []*models.Listing{}
	err := s.DB.Select(&ls, `
		SELECT 
		l.id, l.description, l.price, l.mileage, l.used, l.phone_nos, l.listed_at,
		bike.id "bike.id", bike.model "bike.model", bike.brand "bike.brand", bike.category "bike.category", bike.year "bike.year", bike.engine_capacity "bike.engine_capacity",
		city.id "city.id", city.city "city.city", city.district "city.district", city.province "city.province",
		seller.id "seller.id", seller.first_name "seller.first_name", seller.last_name "seller.last_name", seller.email "seller.email", seller.joined_at "seller.joined_at"
		FROM listings AS l
		INNER JOIN bikes as bike ON l.bike_id = bike.id
		INNER JOIN cities AS city ON l.city_id = city.id
		INNER JOIN users AS seller ON l.seller_id = seller.id
		WHERE l.seller_id = ?
	`, sellerId)

	if err != nil {
		return nil, err
	}

	for _, l := range ls {
		err = s.PopulateListing(l)
		if err != nil {
			return nil, err
		}
	}

	return ls, nil
}

type SortByAttr string

const (
	SORT_BY_PRICE     SortByAttr = "l.price"
	SORT_BY_MILEAGE   SortByAttr = "l.mileage"
	SORT_BY_LISTED_AT SortByAttr = "l.listed_at"
)

type SortInOrder string

const (
	SORT_IN_ASC  SortInOrder = "ASC"
	SORT_IN_DESC SortInOrder = "DESC"
)

type FilterByAttr string

const (
	FILTER_BY_PRICE           FilterByAttr = "l.price"
	FILTER_BY_MILEAGE         FilterByAttr = "l.mileage"
	FILTER_BY_USED            FilterByAttr = "l.used"
	FILTER_BY_BRAND           FilterByAttr = "bike.brand"
	FILTER_BY_MODEL           FilterByAttr = "bike.model"
	FILTER_BY_YEAR            FilterByAttr = "bike.year"
	FILTER_BY_CATEGORY        FilterByAttr = "bike.category"
	FILTER_BY_ENGINE_CAPACITY FilterByAttr = "bike.engine_capacity"
	FILTER_BY_CITY            FilterByAttr = "city.city"
	FILTER_BY_DISTRICT        FilterByAttr = "city.district"
	FILTER_BY_PROVINCE        FilterByAttr = "city.province"
)

type Range struct {
	Min int
	Max int
}

type SearchOptions struct {
	SortBy   SortByAttr
	SortIn   SortInOrder
	Limit    int
	FilterBy map[FilterByAttr]interface{}
}

func (s *ListingService) SearchListings(query string, opts *SearchOptions) ([]*models.Listing, error) {
	q := "%" + strings.Join(strings.Split(query, " "), "%") + "%"

	fq := ""
	fm := map[string]interface{}{
		"query": q,
	}

	if len(opts.FilterBy) > 0 {
		fqs := []string{}

		for attr, val := range opts.FilterBy {
			attr_param := strings.ReplaceAll(string(attr), ".", "_")

			if attr == FILTER_BY_PRICE || attr == FILTER_BY_MILEAGE || attr == FILTER_BY_ENGINE_CAPACITY {
				r, ok := val.(Range)
				if !ok {
					return nil, fmt.Errorf("Invalid filter value")
				}

				if r.Min != -1 {
					fqs = append(fqs, fmt.Sprintf("%s >= :%s_min", attr, attr_param))
					fm[attr_param+"_min"] = strconv.Itoa(r.Min)
				}

				if r.Max != -1 {
					fqs = append(fqs, fmt.Sprintf("%s <= :%s_max", attr, attr_param))
					fm[attr_param+"_max"] = strconv.Itoa(r.Max)
				}
			} else if attr == FILTER_BY_USED {
				b, ok := val.(bool)
				if !ok {
					return nil, fmt.Errorf("Invalid filter value")
				}

				fqs = append(fqs, fmt.Sprintf("%s = :%s", attr, attr_param))
				fm[attr_param] = fmt.Sprintf("%t", b)
			} else {
				s, ok := val.(string)
				if !ok {
					return nil, fmt.Errorf("Invalid filter value")
				}

				fqs = append(fqs, fmt.Sprintf("%s = :%s", attr, attr_param))
				fm[attr_param] = s
			}
		}

		if len(fqs) > 0 {
			fq = "AND " + strings.Join(fqs, " AND ")
		}
	}

	ls := []*models.Listing{}
	stmt, err := s.DB.PrepareNamed(fmt.Sprintf(`
		SELECT 
		l.id, l.description, l.price, l.mileage, l.used, l.phone_nos, l.listed_at,
		bike.id "bike.id", bike.model "bike.model", bike.brand "bike.brand", bike.category "bike.category", bike.year "bike.year", bike.engine_capacity "bike.engine_capacity",
		city.id "city.id", city.city "city.city", city.district "city.district", city.province "city.province",
		seller.id "seller.id", seller.first_name "seller.first_name", seller.last_name "seller.last_name", seller.email "seller.email", seller.joined_at "seller.joined_at"
		FROM listings AS l
		INNER JOIN bikes as bike ON l.bike_id = bike.id
		INNER JOIN cities AS city ON l.city_id = city.id
		INNER JOIN users AS seller ON l.seller_id = seller.id
		WHERE CONCAT(bike.brand, ' ', bike.model, ' ', bike.year, ' ', l.description, ' ', city.city, ' ', city.district, ' ', city.province) LIKE :query
		%s
		ORDER BY %s %s
		LIMIT %d
	`, fq, opts.SortBy, opts.SortIn, opts.Limit))

	fmt.Printf("%+v\n", opts)
	fmt.Printf("%+v\n", fm)
	fmt.Printf("%+v\n", fq)
	//fmt.Printf("%s", stmt.QueryString)

	if err != nil {
		return nil, err
	}

	err = stmt.Select(&ls, fm)

	if err != nil {
		return nil, err
	}

	for _, l := range ls {
		err = s.PopulateListing(l)
		if err != nil {
			return nil, err
		}
	}

	return ls, nil
}

func (s *ListingService) CreateListing(sellerId string, bikeId string, description string, price int, mileage int, used bool, cityId string, phoneNos []string, images []io.Reader) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	_, err = s.DB.Exec("INSERT INTO listings(id, seller_id, bike_id, description, price, mileage, used, city_id, phone_nos) VALUES(?,?,?,?,?,?,?,?,?)", id, sellerId, bikeId, description, price, mileage, used, cityId, strings.Join(phoneNos, ","))
	if err != nil {
		return err
	}

	for _, img := range images {
		imgID, err := uuid.NewRandom()
		if err != nil {
			return err
		}

		err = os.MkdirAll("./assets/dist/listings/"+id.String()+"/", 0755)
		if err != nil {
			return err
		}

		f, err := os.Create("./assets/dist/listings/" + id.String() + "/" + imgID.String())
		if err != nil {
			return err
		}

		_, err = f.ReadFrom(img)
		if err != nil {
			return err
		}

		f.Sync()
		f.Close()
	}

	return nil
}
