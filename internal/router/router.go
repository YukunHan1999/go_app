package router

import (
	atthandler "github.com/myapp/internal/modules/attachment/handler"
	attrepo "github.com/myapp/internal/modules/attachment/repository"
	attsvc "github.com/myapp/internal/modules/attachment/service"

	codetplhandler "github.com/myapp/internal/modules/codetemplate/handler"
	codetplrepo "github.com/myapp/internal/modules/codetemplate/repository"
	codetplsvc "github.com/myapp/internal/modules/codetemplate/service"

	debughandler "github.com/myapp/internal/modules/debug/handler"
	debugrepo "github.com/myapp/internal/modules/debug/repository"
	debugsvc "github.com/myapp/internal/modules/debug/service"

	envziphandler "github.com/myapp/internal/modules/envzip/handler"
	envziprepo "github.com/myapp/internal/modules/envzip/repository"
	envzipsvc "github.com/myapp/internal/modules/envzip/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoute(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// attachmen domain
	attR := attrepo.NewAttachmentRepo(db)
	attS := attsvc.NewAttachmentService(attR)
	attH := atthandler.NewAttachmentHandler(attS)

	// codetemplate domain
	codetplR := codetplrepo.NewCodeTemplateRepo(db)
	codetplS := codetplsvc.NewCodeTemplateService(codetplR)
	codetplH := codetplhandler.NewCodeTemplateHandler(codetplS)

	// debug domain
	directoryR := debugrepo.NewDirectoryRepo(db)
	packageR := debugrepo.NewPackageRepo(db)
	programR := debugrepo.NewProgramRepo(db)
	debuginfoR := debugrepo.NewDebugInfoRepo(db)
	debuginfoS := debugsvc.NewDebugInfoService(debuginfoR, attS)
	directoryS := debugsvc.NewDirectoryService(directoryR)
	programS := debugsvc.NewProgramService(programR, debuginfoS)
	packageS := debugsvc.NewPackageService(db, packageR, programS)
	debugS := debugsvc.NewDebugService(db, debuginfoS, directoryS, packageS, attS)
	debugH := debughandler.NewDirectoryHandler(debugS)

	// envzip domain
	envzipR := envziprepo.NewEnvZipRepo(db)
	envzipS := envzipsvc.NewEnvZipService(envzipR, attS)
	envzipH := envziphandler.NewEnvZipHandler(envzipS)

	// attachment
	r.POST("/upload", attH.UploadFile)
	r.GET("/preview/:filename", attH.PreviewFile)
	r.GET("/clear", attH.ClearGarbageAtt)

	// codetemplate
	r.GET("/code/count", codetplH.Count)
	r.GET("/code", codetplH.FindAll)
	r.POST("/code", codetplH.SaveOrUpdate)
	r.DELETE("/code/:id", codetplH.Delete)

	// debug
	r.POST("/debug/dir", debugH.CreateDir)
	r.PUT("/debug/dir", debugH.UpdateDir)

	r.POST("/debug/pkg", debugH.CreatePkg)
	r.DELETE("/debug/pkg", debugH.ClearPkgInfo)
	r.PUT("/debug/pkg", debugH.UpdatePkg)
	r.GET("/debug/pkg/:pid", debugH.FindPackageById)

	r.DELETE("/debug", debugH.DeletePkgOrDir)
	r.GET("/debug/:pid", debugH.FetchPkgOrDirData)
	r.GET("/debug", debugH.FetchAllDirData)

	// envzip
	r.POST("/env", envzipH.CreateEnv)
	r.DELETE("/env/:id", envzipH.DeleteEnv)
	r.PUT("/env", envzipH.UpdateEnv)
	r.GET("/env/dir", envzipH.FindDir)
	r.GET("/env/:pid", envzipH.FindByPId)
	return r
}
