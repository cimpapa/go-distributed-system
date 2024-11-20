package grades

func init() {
	students = []Student{
		{
			ID: 1,
			FirstName: "lei",
			LastName: "wang",
			Grades: []Grade{
				{
					Title: "FinalExam",
					Type: GradeExam,
					Score: 90.5,
				},
				{
					Title: "ClassQuiz",
					Type: GradeQuiz,
					Score: 88.8,
				},
				{
					Title: "ClassTest",
					Type: GradeTest,
					Score: 97,
				},
			},
		},
		{
			ID: 2,
			FirstName: "luuk",
			LastName: "wang",
			Grades: []Grade{
				{
					Title: "FinalExam",
					Type: GradeExam,
					Score: 70.5,
				},
				{
					Title: "ClassQuiz",
					Type: GradeQuiz,
					Score: 89,
				},
				{
					Title: "ClassTest",
					Type: GradeTest,
					Score: 74,
				},
			},
		},
		{
			ID: 3,
			FirstName: "Cooper",
			LastName: "Guo",
			Grades: []Grade{
				{
					Title: "FinalExam",
					Type: GradeExam,
					Score: 100,
				},
				{
					Title: "ClassQuiz",
					Type: GradeQuiz,
					Score: 51,
				},
				{
					Title: "ClassTest",
					Type: GradeTest,
					Score: 96,
				},
			},
		},
	}
}