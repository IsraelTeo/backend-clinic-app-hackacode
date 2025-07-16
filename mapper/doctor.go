package mapper

import (
	"fmt"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/dto"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
)

// DoctorMapper defines the interface for mapping between Doctor model and DTO
type DoctorMapper interface {
	ModelToResponse(doctor *model.Doctor) (*dto.DoctorResponse, error)
	ResponseToModel(response *dto.DoctorResponse) (*model.Doctor, error)
	ModelToRequest(doctor *model.Doctor) (*dto.DoctorRequest, error)
	RequestToModel(request *dto.DoctorRequest) (*model.Doctor, error)
	ModelListToResponseList(doctors []model.Doctor) ([]dto.DoctorResponse, error)
	ResponseListToRequestList(responses []dto.DoctorResponse) ([]dto.DoctorRequest, error)
}

type doctorMapper struct{}

func NewDoctorMapper() DoctorMapper {
	return &doctorMapper{}
}

// ModelToResponse converts a Doctor model to a DoctorResponse DTO.
func (m *doctorMapper) ModelToResponse(doctor *model.Doctor) (*dto.DoctorResponse, error) {
	if doctor == nil {
		return nil, fmt.Errorf("doctor is nil")
	}
	return &dto.DoctorResponse{
		PersonResponse: dto.PersonResponse{
			ID:          doctor.ID,
			Name:        doctor.Name,
			LastName:    doctor.LastName,
			DNI:         doctor.DNI,
			BirthDate:   doctor.BirthDate,
			Email:       doctor.Email,
			PhoneNumber: doctor.PhoneNumber,
			Address:     doctor.Address,
		},
		Especialty: doctor.Especialty,
		Days:       doctor.Days,
		StartTime:  doctor.StartTime,
		EndTime:    doctor.EndTime,
		Salary:     doctor.Salary,
	}, nil
}

// ResponseToModel converts a DoctorResponse DTO to a Doctor model.
func (m *doctorMapper) ResponseToModel(response *dto.DoctorResponse) (*model.Doctor, error) {
	if response == nil {
		return nil, fmt.Errorf("response is nil")
	}
	return &model.Doctor{
		Person: model.Person{
			Name:        response.Name,
			LastName:    response.LastName,
			DNI:         response.DNI,
			BirthDate:   response.BirthDate,
			Email:       response.Email,
			PhoneNumber: response.PhoneNumber,
			Address:     response.Address,
		},
		Especialty: response.Especialty,
		Days:       response.Days,
		StartTime:  response.StartTime,
		EndTime:    response.EndTime,
		Salary:     response.Salary,
	}, nil
}

// ModelToRequest converts a Doctor model to a DoctorRequest DTO.
func (m *doctorMapper) ModelToRequest(doctor *model.Doctor) (*dto.DoctorRequest, error) {
	if doctor == nil {
		return nil, fmt.Errorf("doctor is nil")
	}
	return &dto.DoctorRequest{
		PersonRequest: dto.PersonRequest{
			Name:        doctor.Name,
			LastName:    doctor.LastName,
			DNI:         doctor.DNI,
			BirthDate:   doctor.BirthDate,
			Email:       doctor.Email,
			PhoneNumber: doctor.PhoneNumber,
			Address:     doctor.Address,
		},
		Especialty: doctor.Especialty,
		Days:       doctor.Days,
		StartTime:  doctor.StartTime,
		EndTime:    doctor.EndTime,
		Salary:     doctor.Salary,
	}, nil
}

// RequestToModel converts a DoctorRequest DTO to a Doctor model.
func (m *doctorMapper) RequestToModel(request *dto.DoctorRequest) (*model.Doctor, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}
	return &model.Doctor{
		Person: model.Person{
			Name:        request.Name,
			LastName:    request.LastName,
			DNI:         request.DNI,
			BirthDate:   request.BirthDate,
			Email:       request.Email,
			PhoneNumber: request.PhoneNumber,
			Address:     request.Address,
		},
		Especialty: request.Especialty,
		Days:       request.Days,
		StartTime:  request.StartTime,
		EndTime:    request.EndTime,
		Salary:     request.Salary,
	}, nil
}

// ModelListToResponseList converts a list of Doctor models to a list of DoctorResponse DTOs.
func (m *doctorMapper) ModelListToResponseList(doctors []model.Doctor) ([]dto.DoctorResponse, error) {
	if len(doctors) == 0 {
		return nil, fmt.Errorf("doctors list is empty")
	}

	responses := make([]dto.DoctorResponse, 0, len(doctors))
	for _, doctor := range doctors {
		dto, err := m.ModelToResponse(&doctor)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *dto)
	}

	return responses, nil
}

// ResponseListToRequestList converts a list of DoctorResponse DTOs to a list of DoctorRequest DTOs.
func (m *doctorMapper) ResponseListToRequestList(responses []dto.DoctorResponse) ([]dto.DoctorRequest, error) {
	if len(responses) == 0 {
		return nil, fmt.Errorf("responses list is empty")
	}

	requests := make([]dto.DoctorRequest, 0, len(responses))
	for _, resp := range responses {
		req := dto.DoctorRequest{
			PersonRequest: dto.PersonRequest{
				Name:        resp.Name,
				LastName:    resp.LastName,
				DNI:         resp.DNI,
				BirthDate:   resp.BirthDate,
				Email:       resp.Email,
				PhoneNumber: resp.PhoneNumber,
				Address:     resp.Address,
			},
			Especialty: resp.Especialty,
			Days:       resp.Days,
			StartTime:  resp.StartTime,
			EndTime:    resp.EndTime,
			Salary:     resp.Salary,
		}
		requests = append(requests, req)
	}

	return requests, nil
}
