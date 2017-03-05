package models

func UsersByPivot(db XODB, pivot, size int64) ([]*User, error) {
	var err error

	const nsqlstr = `SELECT ` +
		`snowflake ` +
		`FROM public.users ` +
		`WHERE snowflake < $1 ` +
		`LIMIT $2`

	const psqlstr = `SELECT ` +
		`snowflake ` +
		`FROM public.users ` +
		`WHERE snowflake > $1 ` +
		`LIMIT $2`

	var sqlstr string
	if size > 0 {
		sqlstr = psqlstr
	} else {
		sqlstr = nsqlstr
		size = size * -1
	}

	XOLog(sqlstr, pivot, size)

	res, err := db.Query(sqlstr, pivot, size)
	if err != nil {
		return nil, err
	}

	var snowflakes = []int64{}

	for res.Next() {
		var sf int64
		err = res.Scan(&sf)
		if err != nil {
			return nil, err
		}
		snowflakes = append(snowflakes, sf)
	}

	users := []*User{}

	for _,v := range snowflakes {
		u, err := UserBySnowflake(db, v)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}