-- https://webapp.io/blog/postgres-is-the-answer/
CREATE OR REPLACE FUNCTION task_queue_status_notify()
RETURNS trigger AS
$$
BEGIN
    PERFORM pg_notify('task_queue_status_channel', NEW.id::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER task_queue_status
    AFTER INSERT OR UPDATE
    OF status ON task_queue
FOR EACH ROW
    EXECUTE PROCEDURE task_queue_status_notify();
