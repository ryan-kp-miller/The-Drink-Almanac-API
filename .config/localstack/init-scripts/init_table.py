import psycopg2
connection = psycopg2.connect(
    host = 'localstack',
    port = 4510,
    user = 'test',
    password = 'test',
    database='task_store'
    )
cursor=connection.cursor()

#creating table passengers
cursor.execute("""CREATE TABLE task_store(
task_id text PRIMARY KEY,
task_reg_ts timestamp,
task_start_ts timestamp,
task_end_ts timestamp,
meta JSONB,
asset_name text,
status text)
""")

connection.commit()