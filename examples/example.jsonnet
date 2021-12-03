{
  simple_spec: spec([
    {
      name: "google homepage",
      request: {
        method: "GET",
        url: "https://www.google.com/"
      },
      expect: {
        code: 200
      }
    }
  ])
}