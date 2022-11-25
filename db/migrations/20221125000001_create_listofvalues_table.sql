create table listofvalues2
(
    listofvalues_type             varchar(50)                  not null,
    system_name                   varchar(255)                 not null,
    display_name                  varchar(255),
    listofvalues_note             varchar(255),
    language_code                 varchar(50)                  not null,
    language_name                 varchar(255)                 not null,
    created_by_date               date    default CURRENT_DATE not null,
    created_by_requestor_username varchar(255)                 not null,
    updated_by_date               date,
    updated_by_requestor_username varchar(255),
    value_is_active               boolean default false        not null,
    parent_listofvalues_type      varchar(50),
    parent_system_name            varchar(255),
    primary key (listofvalues_type, system_name)
);

comment on table listofvalues2 is 'Contains information about list of values';

comment on column listofvalues2.listofvalues_type is 'This is the list of value type';

comment on column listofvalues2.system_name is 'This is the identifier name of the list of value';

comment on column listofvalues2.display_name is 'This is the display name of the list of value';

comment on column listofvalues2.listofvalues_note is 'This is the notes of the list of value';

comment on column listofvalues2.language_code is 'This is the code for the language of the display name value';

comment on column listofvalues2.language_name is 'This is the full name of the language code';

comment on column listofvalues2.created_by_date is 'When the user entry was entered in the system';

comment on column listofvalues2.created_by_requestor_username is 'The Cognito username for the person who is requesting the users creation';

comment on column listofvalues2.updated_by_date is 'When the user entry was last updated in the system';

comment on column listofvalues2.updated_by_requestor_username is 'The Cognito username for the person who is requesting the users update';

comment on column listofvalues2.value_is_active is 'Activates and deactivates a listofvalue';

comment on column listofvalues2.parent_listofvalues_type is 'This is the parent list of value type';

comment on column listofvalues2.parent_system_name is 'This is the identifier name of the list of value parent';