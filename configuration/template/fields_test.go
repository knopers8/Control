package template_test

import (
	"github.com/AliceO2Group/Control/apricot/local"
	"github.com/AliceO2Group/Control/common/gera"
	"github.com/AliceO2Group/Control/configuration"
	"github.com/AliceO2Group/Control/configuration/template"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "text/template"
)

var _ = Describe("Fields", func() {

})

var _ = Describe("Sequence", func() {
	var (
		mockConfSvc configuration.Service
		varStack    template.VarStack
		err         error
	)

	BeforeEach(func() {
		mockConfSvc, err = local.NewService("mock://")
		Expect(err).NotTo(HaveOccurred())
		varStack = template.VarStack{
			Locals:   make(map[string]string),
			Defaults: gera.MakeMap[string, string](),
			Vars:     gera.MakeMap[string, string](),
			UserVars: gera.MakeMap[string, string](),
		}
	})

	Describe("Execute", func() {
		Context("when execution is successful", func() {
			It("should complete without errors", func() {
				sequence := template.Sequence{
					template.STAGE0: template.Fields{&mockField{value: "test"}},
					template.STAGE1: template.Fields{&mockField{value: "test"}},
				}

				//mockConfSvc.(*configuration.MockSource).On("GetComponentConfiguration", componentcfg.Query{}).Return("", nil)

				err := sequence.Execute(
					mockConfSvc,
					"test/path",
					varStack,
					mockBuildObjectStack,
					make(map[string]string),
					make(map[string]Template),
					nil,
					mockStageCallback,
				)

				Expect(err).NotTo(HaveOccurred())
			})
		})
		/*
			Context("when an error occurs during stage execution", func() {
				It("should return the error", func() {
					sequence := template.Sequence{
						template.STAGE0: template.Fields{&mockField{value: "test", err: errors.New("execution error")}},
					}

					//mockConfSvc.(*configuration.MockSource).On("GetComponentConfiguration", componentcfg.Query{}).Return("", nil)

					err := sequence.Execute(
						mockConfSvc,
						"test/path",
						varStack,
						mockBuildObjectStack,
						make(map[string]string),
						make(map[string]Template),
						nil,
						mockStageCallback,
					)

					//Expect(err).To(MatchError("execution error"))
					//Expect(err).To(HaveOccurred())
				})
			})

			Context("when a RoleDisabledError is encountered", func() {
				It("should return the RoleDisabledError", func() {
					sequence := template.Sequence{
						template.STAGE0: template.Fields{&mockField{value: "test", err: &template.RoleDisabledError{RolePath: "test/path"}}},
					}

					//mockConfSvc.(*configuration.MockSource).On("GetComponentConfiguration", componentcfg.Query{}).Return("", nil)

					err := sequence.Execute(
						mockConfSvc,
						"test/path",
						varStack,
						mockBuildObjectStack,
						make(map[string]string),
						make(map[string]Template),
						nil,
						mockStageCallback,
					)

					//Expect(err).To(BeAssignableToTypeOf(&template.RoleDisabledError{}))
					//Expect(err.(*template.RoleDisabledError).RolePath).To(Equal("test/path"))
				})
			}) */
	})
})

// Mock BuildObjectStackFunc
func mockBuildObjectStack(stage template.Stage) map[string]interface{} {
	return map[string]interface{}{"test": "value"}
}

// Mock StageCallbackFunc
func mockStageCallback(stage template.Stage, err error) error {
	return err
}

// Mock Field for testing
type mockField struct {
	value string
	err   error
}

func (m *mockField) Get() string {
	return m.value
}

func (m *mockField) Set(s string) {
	m.value = s
}

//func (m *mockField) Execute(confSvc configuration.Source, parentPath string, varStack map[string]string, objStack map[string]interface{}, baseConfigStack map[string]string, stringTemplateCache map[string]template.Template, workflowRepo repos.IRepo) error {
//	return m.err
//}
