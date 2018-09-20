[
  "./priv/endpoints/endpoints-1.1.csv",
  "./priv/endpoints/endpoints-2.0.csv"
]
|> Enum.each(&IO.inspect(Compliance.Permutations.Generator.process_file(&1)))
