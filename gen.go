package gen

//go:generate mockgen -destination=internal/mock/company_servicer.go -package=mock github.com/djurica-surla/golang-exercise/internal/transport/http CompanyServicer
//go:generate mockgen -destination=internal/mock/token_servicer.go -package=mock github.com/djurica-surla/golang-exercise/internal/transport/http TokenServicer
