package testing

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Model suite
type Model struct {
	suite.Suite
	*require.Assertions
}

// func (m *Model) SetupTest() {
// 	m.Assertions = require.New(m.T())
// }

// // TearDownTest will be called after tests finish
// func (m *Model) TearDownTest() {}

// func NewModel() *Model {
// 	return &Model{}
// }
