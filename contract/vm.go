package contract

import (
	"errors"
	"sync"

	log "github.com/xuperchain/log15"
	"github.com/xuperchain/xuperunion/pb"
	"github.com/xuperchain/xuperunion/xmodel"
)

const (
	// StatusOK is used when contract successfully ends.
	StatusOK = 200
	// StatusErrorThreshold is the status dividing line for the normal operation of the contract
	StatusErrorThreshold = 400
	// StatusError is used when contract fails.
	StatusError = 500
)

var (
	// ErrVMNotExist is returned when found vm not exist
	ErrVMNotExist = errors.New("Vm not exist in vm manager")
)

// Response is the result of the contract run
type Response struct {
	// Status 用于反映合约的运行结果的错误码
	Status int `json:"status"`
	// Message 用于携带一些有用的debug信息
	Message string `json:"message"`
	// Data 字段用于存储合约执行的结果
	Body []byte `json:"body"`
}

func ToPBContractResponse(resp *Response) *pb.ContractResponse {
	return &pb.ContractResponse{
		Status:  int32(resp.Status),
		Message: resp.Message,
		Body:    resp.Body,
	}
}

// ContextConfig define the config of context
type ContextConfig struct {
	XMCache        *xmodel.XMCache
	Initiator      string
	AuthRequire    []string
	ContractName   string
	ResourceLimits Limits
	// Whether contract can be initialized
	CanInitialize bool
}

// VirtualMachine define virtual machine interface
type VirtualMachine interface {
	GetName() string
	NewContext(*ContextConfig) (Context, error)
}

// Context define context interface
type Context interface {
	Invoke(method string, args map[string][]byte) (*Response, error)
	ResourceUsed() Limits
	Release() error
}

// VMManager define VMManager type
type VMManager struct {
	lock   *sync.Mutex
	vms    map[string]VirtualMachine
	logger log.Logger
}

// NewVMManager new an instance of VMManager
func NewVMManager(logger log.Logger) (*VMManager, error) {
	vmMgr := &VMManager{
		lock:   new(sync.Mutex),
		vms:    map[string]VirtualMachine{},
		logger: logger,
	}
	return vmMgr, nil
}

// RegisterVM register an instance of VM into VMManager
func (vmMgr *VMManager) RegisterVM(module string, vm VirtualMachine) error {
	vmMgr.lock.Lock()
	defer vmMgr.lock.Unlock()
	vmMgr.vms[module] = vm
	return nil
}

// GetVM return specific virtual machine instance
func (vmMgr *VMManager) GetVM(module string) (VirtualMachine, error) {
	if vmMgr.vms[module] == nil {
		return nil, ErrVMNotExist
	}
	return vmMgr.vms[module], nil
}
