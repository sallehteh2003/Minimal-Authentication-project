package Handlers

import (
	"github.com/sirupsen/logrus"
	"main/Authentication"
	"main/Database"

	"main/Validation"
)

type Server struct {
	Logger         *logrus.Logger
	Database       *Database.DB
	Validation     *Validation.Validation
	Authentication *Authentication.Authentication
}
