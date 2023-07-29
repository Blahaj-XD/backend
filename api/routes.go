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
			"X-Bank-Authorization",
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
		authGroup.Get("/verify", s.AuthVerify)
	}

	parentAdminGroup := s.app.Group("/parent-admin", middleware.Authenticated)
	{
		parentAdminKidsGroup := parentAdminGroup.Group("/kids")
		{
			parentAdminKidsGroup.Get("/", s.ParentAdminListKids)
			parentAdminKidsGroup.Post("/", s.ParentAdminAddKid)
		}

		parentAdminQuestsGroup := parentAdminGroup.Group("/quests")
		{
			parentAdminQuestsGroup.Post("/", s.ParentAdminCreateQuest)
			parentAdminQuestsGroup.Get("/", s.ParentAdminListQuests)

			parentAdminGroup.Patch("/:id/available",
				s.CreateParentAdminUpdateQuestStatusHandler(repo.QuestStatusAvailable))

			parentAdminGroup.Patch("/:id/ongoing",
				s.CreateParentAdminUpdateQuestStatusHandler(repo.QuestStatusOngoing))

			parentAdminGroup.Patch("/:id/done",
				s.CreateParentAdminUpdateQuestStatusHandler(repo.QuestStatusDone))

			parentAdminGroup.Patch("/:id/approve",
				s.CreateParentAdminUpdateQuestStatusHandler(repo.QuestStatusApproved))

			parentAdminGroup.Patch("/:id/cancel",
				s.CreateParentAdminUpdateQuestStatusHandler(repo.QuestStatusCanceled))
		}

		parentAdminBankGroup := parentAdminGroup.Group("/bank")
		{
			parentAdminBankGroup.Get("/",
				s.ParentAdminBankAccountInfo)

			parentAdminBankGroup.Post("/deposit", s.BankAddBalance)

			parentAdminBankGroup.Get("/transactions",
				s.ParentAdminBankTransactionInfo)

			parentAdminBankGroup.Post("/deposit-kid-account/:kidID",
				s.ParentAdminBankDepositKidAccount)

			parentAdminBankGroup.Post("/withdraw-kid-account/:kidID/approve",
				s.CreateParentAdminDecideKidGoalRequestHandler(repo.KidBalanceRequestStatusApproved))

			parentAdminBankGroup.Post("/withdraw-kid-account/:kidID/reject",
				s.CreateParentAdminDecideKidGoalRequestHandler(repo.KidBalanceRequestStatusRejected))
		}
	}

	kidsGroup := s.app.Group("/kids", middleware.Authenticated)
	{
		kidsDashboardGroup := kidsGroup.Group("/:kidID/dashboard")
		{
			kidsBankGroup := kidsDashboardGroup.Group("/bank")
			{
				kidsBankGroup.Get("/", s.KidDashboardBankAccountInfo)
				kidsBankGroup.Get("/transactions", s.KidDashboardBankTransactionInfo)
				kidsBankGroup.Post("/withdraw", s.KidDashboardBankRequestWithdraw)
			}

			kidsGoalsGroup := kidsDashboardGroup.Group("/goals")
			{
				kidsGoalsGroup.Get("/", s.KidDashboardListGoals)
				kidsGoalsGroup.Post("/", s.KidDashboardCreateGoal)
				kidsGoalsGroup.Post("/:goalID/deposit", s.KidDashboardDepositGoal)

				kidsGoalsGroup.Patch("/:goalID/ongoing",
					s.CreateKidDashboardUpdateGoalStatusHandler(repo.GoalStatusOngoing))

				kidsGoalsGroup.Patch("/:goalID/achieved",
					s.CreateKidDashboardUpdateGoalStatusHandler(repo.GoalStatusAchieved))

				kidsGoalsGroup.Patch("/:goalID/cancelled",
					s.CreateKidDashboardUpdateGoalStatusHandler(repo.GoalStatusCancelled))
			}

			kidsQuestsGroup := kidsDashboardGroup.Group("/quests")
			{
				kidsQuestsGroup.Post("/:questID", s.KidDashboardQuestTake)
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
