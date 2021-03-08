package repositories

import (
	"project/domain/model"
	"project/interfaces"
)

//SampleRepository .
type SampleRepository struct {
	interfaces.IDbHandler
}

//FindSample .
func (r *SampleRepository) FindSample(id int64) (*model.Sample, error) {
	rows, err := r.Query(`
SELECT
	 id
	,code
	,name
FROM
	sample
WHERE
	id =:id
`, map[string]interface{}{"id": id})

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sample := &model.Sample{}
	for rows.Next() {
		err := rows.StructScan(sample)
		if err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sample, nil
}
