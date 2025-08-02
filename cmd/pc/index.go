package pc

import (
	"MrAlbertX/server/cmd/dependencies"

	"github.com/spf13/cobra"
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Escaneia o sistema em busca de programas para criar um cache de busca.",
	Long:  `Executa uma varredura interativa no sistema para encontrar executáveis e atalhos.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		uc := dependencies.GetReindexProgramsUseCase()
		return uc.Execute()
	},
}
