{
  stepsGroups: [
    {
      description: "Default Group",
      required: "all",
      steps: [
        {
          id: "step-1",
          icon: "common",
          name: "Default Step",
          description: "Default Step",
          inputs: [
            {
              id: "mysql-password",
              name: "MySQL password",
              description: "Please entry MySQL password",
              type: "string",
              kind: "hidden",
              rules: "required",
            },
          ],
        },
      ],
    },
  ],
}
