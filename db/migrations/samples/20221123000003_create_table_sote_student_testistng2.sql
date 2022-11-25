create sequence if not exists sote_student_test2_sote_student_test2_id_seq
    increment by 3;
create table if not exists sote_student_test2
(
    sote_student_test2_id bigint default nextval('sote_student_test2_sote_student_test2_id_seq'::regclass) not null
        primary key,
    name                 text                                                                           not null,
    age                 integer                                                                        not null,
    class                text                                                                           not null
);

comment
on column sote_student_test2.sote_student_test2_id is 'This is the unique identifier for each student';

comment
on column sote_student_test2.name is 'This is the name of the student';

comment
on column sote_student_test2.age is 'This is the age of the student';

comment
on column sote_student_test2.class is 'This is the class or grade of the student';
