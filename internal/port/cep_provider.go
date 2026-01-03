package port

import "rfoh/cloud-run/internal/domain/entity"

type CEPProvider interface {
	FindLocationByCEP(cep *entity.CEP) (*entity.Location, error)
}
