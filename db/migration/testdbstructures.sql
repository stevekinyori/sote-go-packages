CREATE SCHEMA sotetest AUTHORIZATION sote;

/*
 * Example of a simple reference table structure.  Note, answertype_note comment is missing as an edge case.
 */
CREATE TABLE sotetest.referencetable
(
    reference_name                VARCHAR(25)                      NOT NULL
        CONSTRAINT referencetable_pkey
            PRIMARY KEY,
    reference_name_display        VARCHAR(250)                     NOT NULL,
    created_by_date               DATE DEFAULT ('now'::TEXT)::DATE NOT NULL,
    created_by_requestor_username VARCHAR(255)                     NOT NULL,
    updated_by_date               DATE,
    updated_by_requestor_username VARCHAR(255)
);

COMMENT ON TABLE sotetest.referencetable IS 'This is a test table and is not to be used except for testing.';

COMMENT ON COLUMN sotetest.referencetable.reference_name IS 'This is the name of the reference type used by Sote System';

COMMENT ON COLUMN sotetest.referencetable.reference_name_display IS 'This is the display name of the reference type used by Sote System';

COMMENT ON COLUMN sotetest.referencetable.created_by_date IS 'When the user entry was entered in the system';

COMMENT ON COLUMN sotetest.referencetable.created_by_requestor_username IS 'The Cognito username for the person who is requesting the users creation';

COMMENT ON COLUMN sotetest.referencetable.updated_by_date IS 'When the user entry was last updated in the system';

COMMENT ON COLUMN sotetest.referencetable.updated_by_requestor_username IS 'The Cognito username for the person who is requesting the users update';

ALTER TABLE sotetest.referencetable
    OWNER TO sote;

INSERT INTO sotetest.referencetable (reference_name, reference_name_display, created_by_date, created_by_requestor_username, updated_by_date, updated_by_requestor_username)
VALUES ('COMPANY_NAME', 'Company Name', '2020-07-27', 'HARVESTER', NULL, NULL);
INSERT INTO sotetest.referencetable (reference_name, reference_name_display, created_by_date, created_by_requestor_username, updated_by_date, updated_by_requestor_username)
VALUES ('PORT_OF_ENTRY', 'Port of Entry', '2020-07-27', 'HARVESTER', '2020-07-27', 'HARVESTER');
INSERT INTO sotetest.referencetable (reference_name, reference_name_display, created_by_date, created_by_requestor_username, updated_by_date, updated_by_requestor_username)
VALUES ('COMPANY_NAME_SEED', 'Company Name Seed', '2020-07-27', 'SEED', NULL, NULL);
INSERT INTO sotetest.referencetable (reference_name, reference_name_display, created_by_date, created_by_requestor_username, updated_by_date, updated_by_requestor_username)
VALUES ('PORT_OF_ENTRY_SEED', 'Port of Entry Seed', '2020-07-27', 'SEED', '2020-07-27', 'HARVESTER');

/*
 * Example of a simple parent table for parent/child structure.
 */
CREATE TABLE sotetest.parenttable
(
    parenttable_id                  BIGSERIAL                           NOT NULL
        CONSTRAINT parenttable_pkey
            PRIMARY KEY,
    row_version                     BIGINT  DEFAULT 1                   NOT NULL,
    user_organizations_id           BIGINT                              NOT NULL,
    user_is_assigned                BOOLEAN DEFAULT FALSE               NOT NULL,
    created_by_date                 DATE    DEFAULT ('now'::TEXT)::DATE NOT NULL,
    created_by_requestor_username   VARCHAR(255)                        NOT NULL,
    updated_by_date                 DATE,
    updated_by_requestor_username   VARCHAR(255),
    user_phone_country_code         VARCHAR(10),
    user_phone_number               VARCHAR(20),
    user_email                      VARCHAR(255),
    user_last_name                  VARCHAR(100),
    user_first_name                 VARCHAR(50),
    user_middle_name                VARCHAR(50),
    user_name_prefix                VARCHAR(5),
    user_name_suffix                VARCHAR(10),
    user_hire_date                  DATE,
    user_driver_license_id          VARCHAR(50),
    user_driver_license_expiry_date DATE,
    user_wage                       REAL,
    user_pay_period                 VARCHAR(15),
    device_token                    VARCHAR(250),
    informational_user              BOOLEAN DEFAULT FALSE               NOT NULL,
    user_is_active                  BOOLEAN DEFAULT FALSE               NOT NULL,
    client_company_id               BIGINT
);

COMMENT ON TABLE sotetest.parenttable IS 'This is a test table and is not to be used except for testing.';

COMMENT ON COLUMN sotetest.parenttable.parenttable_id IS 'This is a system generated id to make entry unique';

COMMENT ON COLUMN sotetest.parenttable.row_version IS 'This is a system generated number that is used to control updates to the record';

COMMENT ON COLUMN sotetest.parenttable.user_organizations_id IS 'This is the organizations the user belongs to and has access too';

COMMENT ON COLUMN sotetest.parenttable.user_is_assigned IS 'User is assigned a task in the system';

COMMENT ON COLUMN sotetest.parenttable.created_by_date IS 'When the user entry was entered in the system';

COMMENT ON COLUMN sotetest.parenttable.created_by_requestor_username IS 'The Cognito username for the person who is requesting the users creation';

COMMENT ON COLUMN sotetest.parenttable.updated_by_date IS 'When the user entry was last updated in the system';

COMMENT ON COLUMN sotetest.parenttable.updated_by_requestor_username IS 'The Cognito username for the person who is requesting the users update';

COMMENT ON COLUMN sotetest.parenttable.user_phone_country_code IS 'User phones country code';

COMMENT ON COLUMN sotetest.parenttable.user_phone_number IS 'User phones number';

COMMENT ON COLUMN sotetest.parenttable.user_email IS 'Users email';

COMMENT ON COLUMN sotetest.parenttable.user_last_name IS 'This is the last name of the user';

COMMENT ON COLUMN sotetest.parenttable.user_first_name IS 'This is the first name of the user';

COMMENT ON COLUMN sotetest.parenttable.user_middle_name IS 'This is the middle name of the user';

COMMENT ON COLUMN sotetest.parenttable.user_name_prefix IS 'The title the user is addressed by';

COMMENT ON COLUMN sotetest.parenttable.user_name_suffix IS 'This is the position, educational degree, accreditation, office, honor of the user';

COMMENT ON COLUMN sotetest.parenttable.user_hire_date IS 'When the user was hired by the organization';

COMMENT ON COLUMN sotetest.parenttable.user_driver_license_id IS 'The identifier on the users license';

COMMENT ON COLUMN sotetest.parenttable.user_driver_license_expiry_date IS 'The date the license expires';

COMMENT ON COLUMN sotetest.parenttable.user_wage IS 'The amount the person is paid for the provided pay period';

COMMENT ON COLUMN sotetest.parenttable.user_pay_period IS 'The frequency the person is paid';

COMMENT ON COLUMN sotetest.parenttable.device_token IS 'The value Google provides to identify a phone';

COMMENT ON COLUMN sotetest.parenttable.informational_user IS 'identifies system user(TRUE) and informational user(FALSE)';

COMMENT ON COLUMN sotetest.parenttable.user_is_active IS 'Activates and deactivates user';

COMMENT ON COLUMN sotetest.parenttable.client_company_id IS 'The system id that is used to control access to organization customer information';

ALTER TABLE sotetest.parenttable
    OWNER TO sote;

CREATE UNIQUE INDEX parenttable_orgid_phone_uindex
    ON sotetest.parenttable (user_phone_number, user_phone_country_code);

CREATE UNIQUE INDEX parenttable_orgid_email_uindex
    ON sotetest.parenttable (user_email);

CREATE INDEX parenttable_orgid_not_uindex
    ON sotetest.parenttable (user_organizations_id);

CREATE INDEX parenttable_clientcompanyid_not_uindex
    ON sotetest.parenttable (client_company_id);

CREATE INDEX parenttable_orgid_clientcompanyid_not_uindex
    ON sotetest.parenttable (client_company_id, user_organizations_id);


/*
 * Example of a simple child table for parent/child structure.
 */
CREATE TABLE sotetest.parentchildtable
(
    parenttable_id   BIGINT NOT NULL,
    cognito_username VARCHAR(255)
);

COMMENT ON TABLE sotetest.parentchildtable IS 'This is a test table and is not to be used except for testing.';

COMMENT ON COLUMN sotetest.parentchildtable.parenttable_id IS 'This is a system generated id to make entry unique';

COMMENT ON COLUMN sotetest.parentchildtable.cognito_username IS 'The Cognito username for this user';

ALTER TABLE sotetest.parentchildtable
    OWNER TO sote;


/*
 * Example of a many to many structure.
 */
CREATE TABLE sotetest.parentsideonetable
(
    parentsideonetable_id         BIGSERIAL                          NOT NULL
        CONSTRAINT parentsideonetable_pkey
            PRIMARY KEY,
    row_version                   BIGINT DEFAULT 1                   NOT NULL,
    created_by_date               DATE   DEFAULT ('now'::TEXT)::DATE NOT NULL,
    created_by_requestor_username VARCHAR(255)                       NOT NULL
);

COMMENT ON TABLE sotetest.parentsideonetable IS 'This contains all vehicle that have been assigned to job cards';

COMMENT ON COLUMN sotetest.parentsideonetable.parentsideonetable_id IS 'This is a system generated id to make entry unique';

COMMENT ON COLUMN sotetest.parentsideonetable.row_version IS 'This is a system generated number that is used to control updates to the record';

COMMENT ON COLUMN sotetest.parentsideonetable.created_by_date IS 'When the user entry was entered in the system';

COMMENT ON COLUMN sotetest.parentsideonetable.created_by_requestor_username IS 'The Cognito username for the person who is requesting the items creation';

ALTER TABLE sotetest.parentsideonetable
    OWNER TO sote;

CREATE TABLE sotetest.parentsidetwotable
(
    parentsidetwotable_id         BIGSERIAL                          NOT NULL
        CONSTRAINT parentsidetwotable_pkey
            PRIMARY KEY,
    row_version                   BIGINT DEFAULT 1                   NOT NULL,
    created_by_date               DATE   DEFAULT ('now'::TEXT)::DATE NOT NULL,
    created_by_requestor_username VARCHAR(255)                       NOT NULL
);

COMMENT ON TABLE sotetest.parentsidetwotable IS 'This contains all vehicle that have been assigned to job cards';

COMMENT ON COLUMN sotetest.parentsidetwotable.parentsidetwotable_id IS 'This is a system generated id to make entry unique';

COMMENT ON COLUMN sotetest.parentsidetwotable.row_version IS 'This is a system generated number that is used to control updates to the record';

COMMENT ON COLUMN sotetest.parentsidetwotable.created_by_date IS 'When the user entry was entered in the system';

COMMENT ON COLUMN sotetest.parentsidetwotable.created_by_requestor_username IS 'The Cognito username for the person who is requesting the items creation';

ALTER TABLE sotetest.parentsidetwotable
    OWNER TO sote;

CREATE TABLE sotetest.manytomanytable
(
    manytomanytable_id            BIGSERIAL                          NOT NULL
        CONSTRAINT manytomanytable_pkey
            PRIMARY KEY,
    parentsideonetable_id         BIGINT                             NOT NULL,
    parentsidetwotable_id         BIGINT                             NOT NULL,
    row_version                   BIGINT DEFAULT 1                   NOT NULL,
    created_by_date               DATE   DEFAULT ('now'::TEXT)::DATE NOT NULL,
    created_by_requestor_username VARCHAR(255)                       NOT NULL
);

COMMENT ON TABLE sotetest.manytomanytable IS 'This contains all vehicle that have been assigned to job cards';

COMMENT ON COLUMN sotetest.manytomanytable.manytomanytable_id IS 'This is a system generated id to make vehicles queue entry unique';

COMMENT ON COLUMN sotetest.parentsidetwotable.row_version IS 'This is a system generated number that is used to control updates to the record';

COMMENT ON COLUMN sotetest.manytomanytable.parentsideonetable_id IS 'This is a system generated id for parentsideonetable';

COMMENT ON COLUMN sotetest.manytomanytable.parentsidetwotable_id IS 'This is a system generated id for parentsidetwotable';

COMMENT ON COLUMN sotetest.manytomanytable.created_by_date IS 'When the user entry was entered in the system';

COMMENT ON COLUMN sotetest.manytomanytable.created_by_requestor_username IS 'The Cognito username for the person who is requesting the items creation';

ALTER TABLE sotetest.manytomanytable
    OWNER TO sote;

CREATE UNIQUE INDEX vehicles_jobcards_uindex
    ON sotetest.manytomanytable (parentsideonetable_id, parentsidetwotable_id);
