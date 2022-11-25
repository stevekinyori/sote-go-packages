# SETUP

Make sure you have run the setup, to set up the necessary migration or seeding files

- You can command line (instructions are in this package's root directory README.md file)
- In code, call sMigration.New function

# NOTES

- Migrations and seeds files can only be .go or .sql types
- .sql types should have valid SQL queries and .go types should have valid GOLANG code referencing sDatabase
  methods|functions
- Each file should have a unique 14 digit prefix(preferably timestamp YMD like `20221124000001`) followed by _ then
  alphanumeric|alpha characters separated by _ (Numerics only are invalid)
- For .go file, the alphanumeric|alpha characters separated by _, forms the method name of that class.eg
  `20221124000001_create_table_students.go` the method name is CreateTableStudents
- The migration and seeding happens in ascending order based on the unique 14 digit prefix on the file names
- Examples of the structure of these files in found in db folder in this package's root directory
- For more information, read the comments in the code