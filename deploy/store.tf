resource "google_bigquery_dataset" "disco" {
  dataset_id                  = "disco"
  friendly_name               = "disco"
  description                 = "disco dataset"
  location                    = "US"
  
  labels = {
    env = "demo"
  }
}

resource "google_bigquery_table" "license" {
  dataset_id = google_bigquery_dataset.disco.dataset_id
  table_id   = "licenses"

  time_partitioning {
    type = "MONTH"
  }

  labels = {
    env = "demo"
  }

  schema = <<EOF
[
    {
        "name": "batch_id",
        "type": "INTEGER",
        "mode": "REQUIRED"
    },
    {
        "name": "image",
        "type": "STRING",
        "mode": "REQUIRED"
    },
    {
        "name": "sha",
        "type": "STRING",
        "mode": "NULLABLE"
    },
    {
        "name": "name",
        "type": "STRING",
        "mode": "REQUIRED"
    },
    {
        "name": "package",
        "type": "STRING",
        "mode": "NULLABLE"
    },
    {
        "name": "updated",
        "type": "TIMESTAMP",
        "mode": "REQUIRED"
    }
]
EOF

}


resource "google_bigquery_table" "vulnerability" {
  dataset_id = google_bigquery_dataset.disco.dataset_id
  table_id   = "vulnerabilities"

  time_partitioning {
    type = "MONTH"
  }

  labels = {
    env = "demo"
  }

  schema = <<EOF
[
    {
        "name": "batch_id",
        "type": "INTEGER",
        "mode": "REQUIRED"
    },
    {
        "name": "image",
        "type": "STRING",
        "mode": "REQUIRED"
    },
    {
        "name": "sha",
        "type": "STRING",
        "mode": "NULLABLE"
    },
    {
        "name": "cve",
        "type": "STRING",
        "mode": "REQUIRED"
    },
    {
        "name": "severity",
        "type": "STRING",
        "mode": "NULLABLE"
    },
    {
        "name": "package",
        "type": "STRING",
        "mode": "NULLABLE"
    },
    {
        "name": "version",
        "type": "STRING",
        "mode": "NULLABLE"
    },
    {
        "name": "title",
        "type": "STRING",
        "mode": "NULLABLE"
    },
    {
        "name": "description",
        "type": "STRING",
        "mode": "NULLABLE"
    },
    {
        "name": "url",
        "type": "STRING",
        "mode": "NULLABLE"
    },
    {
        "name": "updated",
        "type": "TIMESTAMP",
        "mode": "REQUIRED"
    }
]
EOF

}



resource "google_bigquery_table" "package" {
  dataset_id = google_bigquery_dataset.disco.dataset_id
  table_id   = "packages"

  time_partitioning {
    type = "MONTH"
  }

  labels = {
    env = "demo"
  }

  schema = <<EOF
[
    {
        "name": "batch_id",
        "type": "INTEGER",
        "mode": "REQUIRED"
    },
    {
        "name": "image",
        "type": "STRING",
        "mode": "REQUIRED"
    },
    {
        "name": "sha",
        "type": "STRING",
        "mode": "NULLABLE"
    },
    {
        "name": "format",
        "type": "STRING",
        "mode": "REQUIRED"
    },
    {
        "name": "provider",
        "type": "STRING",
        "mode": "REQUIRED"
    },
    {
        "name": "originator",
        "type": "STRING",
        "mode": "NULLABLE"
    },
    {
        "name": "package",
        "type": "STRING",
        "mode": "REQUIRED"
    },
    {
        "name": "version",
        "type": "STRING",
        "mode": "REQUIRED"
    },
    {
        "name": "source",
        "type": "STRING",
        "mode": "NULLABLE"
    },
    {
        "name": "license",
        "type": "STRING",
        "mode": "NULLABLE"
    },
    {
        "name": "updated",
        "type": "TIMESTAMP",
        "mode": "REQUIRED"
    }
]
EOF

}