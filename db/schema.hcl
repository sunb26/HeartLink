table "app_login" {
  schema  = schema.public
  comment = "This table contains the login credentials of each patient for the phone app."
  column "patient_id" {
    null = false
    type = integer
  }
  column "username" {
    null = false
    type = text
  }
  column "password" {
    null = false
    type = text
  }
  primary_key {
    columns = [column.patient_id]
  }
  foreign_key "patient_id" {
    columns     = [column.patient_id]
    ref_columns = [table.patient.column.patient_id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
}
table "patient" {
  schema  = schema.public
  comment = "This table contains all attributes that represent a patient."
  column "physician_id" {
    null = false
    type = text
  }
  column "patient_id" {
    null = false
    type = serial
  }
  column "firstname" {
    null = false
    type = text
  }
  column "lastname" {
    null = false
    type = text
  }
  column "email" {
    null = true
    type = text
  }
  column "height" {
    null = true
    type = integer
  }
  column "weight" {
    null = true
    type = integer
  }
  column "sex" {
    null = true
    type = enum.sex
  }
  column "dob" {
    null    = false
    type    = date
    comment = "Date of Birth"
  }
  primary_key {
    columns = [column.patient_id]
  }
  foreign_key "physician_id" {
    columns     = [column.physician_id]
    ref_columns = [table.physician.column.physician_id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
}
table "physician" {
  schema  = schema.public
  comment = "This table contains all attributes that represent a physician."
  column "physician_id" {
    null    = false
    type    = text
    comment = "This id is varchar instead of serial because it uses the Clerk generated ID."
  }
  column "email" {
    null = false
    type = text
  }
  primary_key {
    columns = [column.physician_id]
  }
}
table "recordings" {
  schema  = schema.public
  comment = "This table tracks the metadata associated with each heart sound recording."
  column "patient_id" {
    null = false
    type = integer
  }
  column "recording_id" {
    null = false
    type = serial
  }
  column "recording_datetime" {
    null = false
    type = timestamptz
  }
  column "download_url" {
    null = false
    type = text
  }
  column "status" {
    null    = false
    type    = enum.view_status
    default = "notSubmitted"
  }
  primary_key {
    columns = [column.recording_id]
  }
  foreign_key "patient_id" {
    columns     = [column.patient_id]
    ref_columns = [table.patient.column.patient_id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
}
enum "sex" {
  schema = schema.public
  values = ["M", "F", "NB"]
}
enum "view_status" {
  schema = schema.public
  values = ["notSubmitted", "pending", "viewed"]
}
schema "public" {
  comment = "standard public schema"
}
