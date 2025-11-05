package mapreduce

import (
	"fmt"
	"net"
	"net/rpc"
)

type Application struct {
	Master        *Master
	AppAddress    string
	masterStarted bool
	listener      net.Listener
}

type RunTaskArgs struct {
	JobName    string
	Files      []string
	MapName    string
	ReduceName string
}
type RunTaskReply struct {
	Accepted bool
	Message  string
}

func NewApplication(appAddress string, MasterAddress string) (application *Application) {

	application = new(Application)
	application.Master = newMaster(MasterAddress)
	application.AppAddress = appAddress
	application.Master.startRPCServer()
	application.masterStarted = true
	return application
}

// Listen starts the Application’s own RPC server so users can call RunTask.
func (app *Application) Listen() error {
	if err := rpc.RegisterName("Application", app); err != nil {
		return fmt.Errorf("register Application RPC: %w", err)
	}
	l, err := net.Listen("tcp", app.AppAddress)
	if err != nil {
		return fmt.Errorf("listen %s: %w", app.AppAddress, err)
	}
	app.listener = l
	go rpc.Accept(l)
	return nil
}

// RunTask kicks off the job. Signature follows net/rpc rules: (args, reply) error.
func (app *Application) RunTask(args *RunTaskArgs, reply *RunTaskReply) error {
	if args == nil {
		*reply = RunTaskReply{Accepted: false, Message: "nil args"}
		return nil
	}

	// Decide which parameters to use: prefer args if provided, else defaults.
	jobName := args.JobName
	files := args.Files
	mapName := args.MapName
	reduceName := args.ReduceName

	// Store the function identifiers in the Master so workers can fetch them.
	app.Master.SetFunctionSpec(FunctionSpec{
		MapName:    mapName,
		ReduceName: reduceName,
	})

	// Start the job asynchronously.
	go app.Master.run(jobName, files, app.Master.schedule, func() {
		// app.Master.stats = app.Master.killWorkers()
		// app.Master.stopRPCServer()
	})

	*reply = RunTaskReply{Accepted: true, Message: "job started"}
	return nil
}

// Close the app RPC server (and optionally the master’s).
func (app *Application) Close() error {
	if app.listener != nil {
		_ = app.listener.Close()
	}
	// You can stop the master RPC in your finish() callback inside run(),
	// so no need to stop it here—unless you want to tear everything down now.
	return nil
}
