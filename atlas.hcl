data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./scripts/atlas-gorm-loader.go",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url

  migration {
    dir = "file://migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
