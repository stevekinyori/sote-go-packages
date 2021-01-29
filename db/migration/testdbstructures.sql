DROP SCHEMA sotetest CASCADE;

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
VALUES ('TEST_COMPANY_NAME', 'TEST Company Name', '2020-07-27', 'HARVESTER', NULL, NULL);
INSERT INTO sotetest.referencetable (reference_name, reference_name_display, created_by_date, created_by_requestor_username, updated_by_date, updated_by_requestor_username)
VALUES ('TEST_PORT_OF_ENTRY', 'TEST Port of Entry', '2020-07-27', 'HARVESTER', '2020-07-27', 'HARVESTER');
INSERT INTO sotetest.referencetable (reference_name, reference_name_display, created_by_date, created_by_requestor_username, updated_by_date, updated_by_requestor_username)
VALUES ('TEST_COMPANY_NAME_SEED', 'TEST Company Name Seed', '2020-07-27', 'SEED', NULL, NULL);
INSERT INTO sotetest.referencetable (reference_name, reference_name_display, created_by_date, created_by_requestor_username, updated_by_date, updated_by_requestor_username)
VALUES ('TEST_PORT_OF_ENTRY_SEED', 'TEST Port of Entry Seed', '2020-07-27', 'SEED', '2020-07-27', 'HARVESTER');

/*
 * Example of a simple parent table for parent/child structure.
 */
CREATE TABLE sotetest.parenttable
(
    parenttable_id                BIGSERIAL                           NOT NULL
        CONSTRAINT parenttable_pkey
            PRIMARY KEY,
    row_version                   BIGINT  DEFAULT 1                   NOT NULL,
    parent_table_name             VARCHAR(100)                        NOT NULL,
    boolean_default_false         BOOLEAN DEFAULT FALSE               NOT NULL,
    created_by_date               DATE    DEFAULT ('now'::TEXT)::DATE NOT NULL,
    created_by_requestor_username VARCHAR(255)                        NOT NULL,
    updated_by_date               DATE,
    updated_by_requestor_username VARCHAR(255),
    general_string_nullable       VARCHAR(255),
    date_no_default               DATE,
    float_no_default              REAL
);

COMMENT ON TABLE sotetest.parenttable IS 'This is a test table and is not to be used except for testing.';

COMMENT ON COLUMN sotetest.parenttable.parenttable_id IS 'This is a system generated id to make entry unique';

COMMENT ON COLUMN sotetest.parenttable.row_version IS 'This is a system generated number that is used to control updates to the record';

COMMENT ON COLUMN sotetest.parenttable.parent_table_name IS 'Test name used for the parent table row';

COMMENT ON COLUMN sotetest.parenttable.boolean_default_false IS 'Tests a boolean column with default false';

COMMENT ON COLUMN sotetest.parenttable.created_by_date IS 'When the user entry was entered in the system';

COMMENT ON COLUMN sotetest.parenttable.created_by_requestor_username IS 'The Cognito username for the person who is requesting the users creation';

COMMENT ON COLUMN sotetest.parenttable.updated_by_date IS 'When the user entry was last updated in the system';

COMMENT ON COLUMN sotetest.parenttable.updated_by_requestor_username IS 'The Cognito username for the person who is requesting the users update';

COMMENT ON COLUMN sotetest.parenttable.general_string_nullable IS 'General string column that is nullable';

COMMENT ON COLUMN sotetest.parenttable.date_no_default IS 'Date column that is nullable';

COMMENT ON COLUMN sotetest.parenttable.float_no_default IS 'Date column that is nullable';

ALTER TABLE sotetest.parenttable
    OWNER TO sote;

CREATE UNIQUE INDEX parenttable_general_string_uindex
    ON sotetest.parenttable (parent_table_name);

CREATE INDEX parenttable_orgid_not_uindex
    ON sotetest.parenttable (float_no_default);


CREATE INDEX parenttable_orgid_clientcompanyid_not_uindex
    ON sotetest.parenttable (boolean_default_false, float_no_default);

INSERT INTO sotetest.parenttable (parenttable_id, row_version, parent_table_name, boolean_default_false, created_by_date, created_by_requestor_username, updated_by_date, updated_by_requestor_username,
                                  general_string_nullable, date_no_default, float_no_default)
VALUES (12453, 1, 'name 1', FALSE, '2020-08-25', 'tyfg4581-jiuh2-3325jghrfd', '2020-08-25', 'tyfg4581-jiuh2-3325jghrfd', 'test field 1', NULL, 12.26);
INSERT INTO sotetest.parenttable (parenttable_id, row_version, parent_table_name, boolean_default_false, created_by_date, created_by_requestor_username, updated_by_date, updated_by_requestor_username,
                                  general_string_nullable, date_no_default, float_no_default)
VALUES (2, 1, 'name 2', FALSE, '2020-08-25', 'seed', NULL, NULL, 'test field 2', NULL, NULL);
INSERT INTO sotetest.parenttable (parenttable_id, row_version, parent_table_name, boolean_default_false, created_by_date, created_by_requestor_username, updated_by_date, updated_by_requestor_username,
                                  general_string_nullable, date_no_default, float_no_default)
VALUES (54201, 1, 'name 3', FALSE, '2020-08-25', 'tyfg4581-jiuh2-3325jghrfd', NULL, NULL, 'test field 5', '2020-08-04', NULL);
INSERT INTO sotetest.parenttable (parenttable_id, row_version, parent_table_name, boolean_default_false, created_by_date, created_by_requestor_username, updated_by_date, updated_by_requestor_username,
                                  general_string_nullable, date_no_default, float_no_default)
VALUES (3, 1, 'name 4', TRUE, '2020-08-25', 'seed', '2020-08-13', 'tyfg4581-jiuh2-3325jghrfd', 'test fld 7', '2020-06-01', 1253.241);

/*
 * Example of a simple child table for parent/child structure.
 */
CREATE TABLE sotetest.parentchildtable
(
    childtable_id    BIGSERIAL NOT NULL
        CONSTRAINT childtable_pkey
            PRIMARY KEY,
    parenttable_id   BIGINT    NOT NULL,
    cognito_username VARCHAR(255),
    reference_name   VARCHAR(25)
);

COMMENT ON TABLE sotetest.parentchildtable IS 'This is a test table and is not to be used except for testing.';

COMMENT ON COLUMN sotetest.parentchildtable.parenttable_id IS 'This is a system generated id to make entry unique';

COMMENT ON COLUMN sotetest.parentchildtable.cognito_username IS 'The Cognito username for this user';

COMMENT ON COLUMN sotetest.parentchildtable.reference_name IS 'This is the name of the reference type used by Sote System';

ALTER TABLE sotetest.parentchildtable
    OWNER TO sote;

CREATE UNIQUE INDEX childtable_parenttable_uindex
    ON sotetest.parenttable (parenttable_id);

INSERT INTO sotetest.parentchildtable (childtable_id, parenttable_id, cognito_username, reference_name)
VALUES (1, 12345, 'Scott', 'TEST_PORT_OF_ENTRY');
INSERT INTO sotetest.parentchildtable (childtable_id, parenttable_id, cognito_username, reference_name)
VALUES (2, 12345, 'April', NULL);
INSERT INTO sotetest.parentchildtable (childtable_id, parenttable_id, cognito_username, reference_name)
VALUES (5, 12345, 'Charity', NULL);
INSERT INTO sotetest.parentchildtable (childtable_id, parenttable_id, cognito_username, reference_name)
VALUES (3, 12345, 'Pam', 'TEST_COMPANY_NAME_SEED');
INSERT INTO sotetest.parentchildtable (childtable_id, parenttable_id, cognito_username, reference_name)
VALUES (4, 12345, 'Steve', 'TEST_COMPANY_NAME_SEED');
INSERT INTO sotetest.parentchildtable (childtable_id, parenttable_id, cognito_username, reference_name)
VALUES (6, 2, NULL, NULL);
INSERT INTO sotetest.parentchildtable (childtable_id, parenttable_id, cognito_username, reference_name)
VALUES (7, 2, NULL, NULL);
INSERT INTO sotetest.parentchildtable (childtable_id, parenttable_id, cognito_username, reference_name)
VALUES (8, 2, NULL, NULL);
INSERT INTO sotetest.parentchildtable (childtable_id, parenttable_id, cognito_username, reference_name)
VALUES (9, 2, NULL, NULL);
INSERT INTO sotetest.parentchildtable (childtable_id, parenttable_id, cognito_username, reference_name)
VALUES (10, 2, NULL, NULL);

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

INSERT INTO sotetest.parentsideonetable (parentsideonetable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (100, 1, '2020-08-25', 'asdf-wert-sdfg-465');
INSERT INTO sotetest.parentsideonetable (parentsideonetable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (200, 1, '2020-08-25', 'seed');
INSERT INTO sotetest.parentsideonetable (parentsideonetable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (300, 1, '2020-08-25', 'asdf-qwerdfgh-dfghdfgh');
INSERT INTO sotetest.parentsideonetable (parentsideonetable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (400, 1, '2020-08-25', 'asdf-qwerdfgh-dfghdfgh');
INSERT INTO sotetest.parentsideonetable (parentsideonetable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (500, 1, '2020-08-25', 'asdf-qwerdfgh-dfghdfgh');
INSERT INTO sotetest.parentsideonetable (parentsideonetable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (600, 1, '2020-08-25', 'asdf-wert-sdfg-465');
INSERT INTO sotetest.parentsideonetable (parentsideonetable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (700, 1, '2020-08-25', 'asdf-wert-sdfg-465');

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

INSERT INTO sotetest.parentsidetwotable (parentsidetwotable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (100100, 1, '2020-08-25', 'asdf-wert-sdfg-465');
INSERT INTO sotetest.parentsidetwotable (parentsidetwotable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (100200, 1, '2020-08-25', 'seed');
INSERT INTO sotetest.parentsidetwotable (parentsidetwotable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (100300, 1, '2020-08-25', 'asdf-qwerdfgh-dfghdfgh');
INSERT INTO sotetest.parentsidetwotable (parentsidetwotable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (100400, 1, '2020-08-25', 'asdf-qwerdfgh-dfghdfgh');
INSERT INTO sotetest.parentsidetwotable (parentsidetwotable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (100500, 1, '2020-08-25', 'asdf-qwerdfgh-dfghdfgh');
INSERT INTO sotetest.parentsidetwotable (parentsidetwotable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (100600, 1, '2020-08-25', 'asdf-wert-sdfg-465');
INSERT INTO sotetest.parentsidetwotable (parentsidetwotable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (100700, 1, '2020-08-25', 'asdf-wert-sdfg-465');

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

INSERT INTO sotetest.manytomanytable (manytomanytable_id, parentsideonetable_id, parentsidetwotable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (1, 100, 100300, 1, '2020-08-25', 'seed');
INSERT INTO sotetest.manytomanytable (manytomanytable_id, parentsideonetable_id, parentsidetwotable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (2, 300, 100300, 1, '2020-08-25', 'asdf-werty-dfgh-655465');
INSERT INTO sotetest.manytomanytable (manytomanytable_id, parentsideonetable_id, parentsidetwotable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (3, 200, 100100, 1, '2020-08-25', 'asdf-werty-dfgh-655465');
INSERT INTO sotetest.manytomanytable (manytomanytable_id, parentsideonetable_id, parentsidetwotable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (4, 200, 100200, 1, '2020-08-25', 'asdf-werty-dfgh-655465');
INSERT INTO sotetest.manytomanytable (manytomanytable_id, parentsideonetable_id, parentsidetwotable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (5, 400, 100600, 1, '2020-08-25', 'seed');
INSERT INTO sotetest.manytomanytable (manytomanytable_id, parentsideonetable_id, parentsidetwotable_id, row_version, created_by_date, created_by_requestor_username)
VALUES (6, 100, 100100, 1, '2020-08-25', 'asdf-werty-dfgh-655465');