package api

import (
	"strings"

	"github.com/BlahajXD/backend/api/middleware"
	"github.com/BlahajXD/backend/repo"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func (s *Server) SetupRoutes() {
	corsMiddleware := cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPatch,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodHead,
			fiber.MethodOptions,
		}, ","),
		AllowHeaders: strings.Join([]string{
			fiber.HeaderAuthorization,
			fiber.HeaderContentType,
			fiber.HeaderAccept,
		}, ","),
	})

	s.app.Use(corsMiddleware)
	s.app.Use(helmet.New())
	s.app.Use(recover.New())

	s.app.Get("/", s.Health)

	authGroup := s.app.Group("/auth")
	{
		authGroup.Post("/register", s.AuthRegister)
		authGroup.Post("/login", s.Login)
	}

	parentAdminGroup := s.app.Group("/parent-admin", middleware.Authenticated)
	{
		parentAdminGroup.Get("/kids", s.ParentAdminListKids)
		parentAdminGroup.Post("/kids", s.ParentAdminAddKid)
		parentAdminGroup.Post("/quests", s.ParentAdminCreateQuest)
		parentAdminGroup.Get("/quests", s.ParentAdminListQuests)

		parentAdminGroup.Patch("/quests/:id/available",
			s.CreateParentAdminUpdateQuestStatusHandler(repo.QuestStatusAvailable))

		parentAdminGroup.Patch("/quests/:id/ongoing",
			s.CreateParentAdminUpdateQuestStatusHandler(repo.QuestStatusOngoing))

		parentAdminGroup.Patch("/quests/:id/done",
			s.CreateParentAdminUpdateQuestStatusHandler(repo.QuestStatusDone))

		parentAdminGroup.Patch("/quests/:id/approve",
			s.CreateParentAdminUpdateQuestStatusHandler(repo.QuestStatusApproved))

		parentAdminGroup.Patch("/quests/:id/cancel",
			s.CreateParentAdminUpdateQuestStatusHandler(repo.QuestStatusCanceled))

		parentAdminGroup.Get("/bank",
			s.ParentAdminBankAccountInfo)

		parentAdminGroup.Get("/bank/transactions",
			s.ParentAdminBankTransactionInfo)

		parentAdminGroup.Post("/bank/deposit-kid-account/:kidID",
			s.ParentAdminBankDepositKidAccount)
	}

	bankGroup := s.app.Group("/bank", middleware.Authenticated)
	{
		bankGroup.Post("/add-balance", s.BankAddBalance)
	}

	kidsGroup := s.app.Group("/kids", middleware.Authenticated)
	{
		kidsDashboardGroup := kidsGroup.Group("/:kidID/dashboard")
		{
			kidsBankGroup := kidsDashboardGroup.Group("/bank")
			{
				kidsBankGroup.Get("/", s.KidDashboardBankAccountInfo)
				kidsBankGroup.Get("/transactions", s.KidDashboardBankTransactionInfo)
			}

			kidsGoalsGroup := kidsDashboardGroup.Group("/goals")
			{
				kidsGoalsGroup.Get("/", s.KidDashboardListGoals)
				kidsGoalsGroup.Post("/", s.KidDashboardCreateGoal)
				kidsGoalsGroup.Post("/:goalID/deposit", s.KidDashboardDepositGoal)
			}
		}
		// kidsGroup.Get("/", srv.GetKidProfile)

		// kidsGroup.Get("/goals", srv.ListKidGoals)
		// kidsGroup.Post("/goals", srv.AddKidGoal)
		// kidsGroup.Get("/goals/:id", srv.GetKidGoal)
		// kidsGroup.Delete("/goals/:id", srv.DeleteKidGoal)

		// kidsGroup.Get("/quests", srv.ListAvailableQuests)
		// kidsGroup.Patch("/quests/take", srv.TakeQuest)
		// kidsGroup.Patch("/quests/complete", srv.CompleteQuest)
	}
}
