package config

// Registry is the configuration necessary to access and interact with a container registry
type Registry struct {
	Username       string `required:"true" envconfig:"USERNAME"`
	Password       string `required:"true" envconfig:"PASSWORD"`
	PrismaURL      string `required:"true" envconfig:"PRISMA_URL"`
	PrismaUsername string `required:"true" envconfig:"PRISMA_USERNAME"`
	PrismaPassword string `required:"true" envconfig:"PRISMA_PASSWORD"`
}
