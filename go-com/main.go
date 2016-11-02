// 
//  2016 October 17
//  John Gilliland [john.gilliland@rndgroup.com]
//  

/*
    Package main is the entry point for the RPC Server that hosts the 
    instrument provider API.
*/
package main


import (    
    	"fmt"
	"log"
	"syscall"
	"unsafe"
	"strconv"

	"github.com/elusive/instrument-api/contracts"
    "github.com/elusive/instrument-api/server/svclog"
    "v.io/x/lib/vlog"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

const progid string  = "HamiltonContext.Executor"

// EventReceiver type is defined to hold and event sink
type EventReceiver struct {
	lpVtbl *EventReceiverVtbl
	ref    int32
	host   *ole.IDispatch
}

// EventReceiverVtbl type is a struct for the event sink virtual table
type EventReceiverVtbl struct {
	pQueryInterface   uintptr
	pAddRef           uintptr
	pRelease          uintptr
	pGetTypeInfoCount uintptr
	pGetTypeInfo      uintptr
	pGetIDsOfNames    uintptr
	pInvoke           uintptr
}

// QueryInterface method is implementation that sets the 
func QueryInterface(this *ole.IUnknown, iid *ole.GUID, punk **ole.IUnknown) uint32 {
	*punk = nil
	if ole.IsEqualGUID(iid, ole.IID_IUnknown) ||
		ole.IsEqualGUID(iid, ole.IID_IDispatch) {
		AddRef(this)
		*punk = this
		return ole.S_OK
	}
	return ole.E_NOINTERFACE
}

// AddRef function increments the reference count
func AddRef(this *ole.IUnknown) int32 {
	pthis := (*EventReceiver)(unsafe.Pointer(this))
	pthis.ref++
	return pthis.ref
}

// Release function decrements the reference count 
func Release(this *ole.IUnknown) int32 {
	pthis := (*EventReceiver)(unsafe.Pointer(this))
	pthis.ref--
	return pthis.ref
}

// GetIDsOfNames function populates array of names with id numbers
func GetIDsOfNames(this *ole.IUnknown, iid *ole.GUID, wnames []*uint16, namelen int, lcid int, pdisp []int32) uintptr {
	for n := 0; n < namelen; n++ {
		pdisp[n] = int32(n)
	}
	return uintptr(ole.S_OK)
}

// GetTypeInfoCount function returns the current type count
func GetTypeInfoCount(pcount *int) uintptr {
	if pcount != nil {
		*pcount = 0
	}
	return uintptr(ole.S_OK)
}

// GetTypeInfo method returns a not implemented type identifier
func GetTypeInfo(ptypeif *uintptr) uintptr {
	return uintptr(ole.E_NOTIMPL)
}

// Invoke method is where the dispatches are handled by number
func Invoke(this *ole.IDispatch, dispid int, riid *ole.GUID, lcid int, flags int16, dispparams *ole.DISPPARAMS, result *ole.VARIANT, pexcepinfo *ole.EXCEPINFO, nerr *uint) uintptr {
	switch dispid {
	case 1:
		vlog.Info("Event received for status change: ")
		executor := (*EventReceiver)(unsafe.Pointer(this)).host
		methodStateEnum, _ := oleutil.CallMethod(executor, "GetMethodState")
		methodStateString := methodStateEnum.ToString()
		methodStateInt, _ := strconv.Atoi(methodStateString)
		vlog.Info("new method state: ", contracts.MethodState(methodStateInt).String())
	default:
		log.Println(dispid)
	}
	return ole.E_NOTIMPL
}


// logs msg if there was an error
func checkForError(err error, msg string) {
	if err != nil {
		vlog.Info(msg)
	}
}

// initialize ole instance of hamilton executor COM object
func initialize() *ole.IDispatch {
	unknown, instError := oleutil.CreateObject(progid)
	checkForError(instError, "Could not create OLE connection to Hamilton Executor.")
	executor, execError := unknown.QueryInterface(ole.IID_IDispatch)
	checkForError(execError, "Could not start a Hamlton Executor instance")
	return executor
}

func out(str string) {
	fmt.Println(str)
}


//
//	MAIN
//
func hxexecute() {
	ole.CoInitialize(0)
	executor := initialize()
	var op = "ENRICHMENT:capture2"
	oleutil.MustCallMethod(executor, "SetOperation", op)
	operation, err := oleutil.GetProperty(executor, "Operation")
	if err != nil {
		fmt.Printf(err.Error())
	}
	fmt.Println("operation:", operation.Val)
}


func main() {
    svclog.Start()

    vlog.Info("Starting GO2COM layer...")

	ole.CoInitialize(0)
	executor := initialize()
	iid, _ := oleutil.ClassIDFrom(progid)

	dest := &EventReceiver{}
	dest.lpVtbl = &EventReceiverVtbl{}
	dest.lpVtbl.pQueryInterface = syscall.NewCallback(QueryInterface)
	dest.lpVtbl.pAddRef = syscall.NewCallback(AddRef)
	dest.lpVtbl.pRelease = syscall.NewCallback(Release)
	dest.lpVtbl.pGetTypeInfoCount = syscall.NewCallback(GetTypeInfoCount)
	dest.lpVtbl.pGetTypeInfo = syscall.NewCallback(GetTypeInfo)
	//dest.lpVtbl.pGetIDsOfNames = syscall.NewCallback(GetIDsOfNames)
	dest.lpVtbl.pInvoke = syscall.NewCallback(Invoke)
	dest.host = executor

	oleutil.ConnectObject(executor, iid, (*ole.IUnknown)(unsafe.Pointer(dest)))
	vlog.Info("Running an med file to check for state change...")
	oleutil.CallMethod(executor, "Run", "C:\\Program Files (x86)\\HAMILTON\\Methods\\Grail\\EnrichmentDriver.hsl")

	// ? how to wait here?
	var m ole.Msg
	for dest.ref != 0 {
		ole.GetMessage(&m, 0, 0, 0)
		ole.DispatchMessage(&m)
	}
}

