package web

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/utils/stringutils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/gin-gonic/gin"
)

// EVMForwardersController manages EVM chains.
type EVMForwardersController struct {
	App chainlink.Application
}

// Index lists EVM chains.
func (cc *EVMForwardersController) Index(c *gin.Context, size, page, offset int) {
	orm := forwarders.NewORM(cc.App.GetSqlxDB(), cc.App.GetLogger(), cc.App.GetConfig())
	fwds, count, err := orm.FindForwarders(0, size)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	var resources []presenters.EVMForwarderResource
	for _, fwd := range fwds {
		resources = append(resources, presenters.NewEVMForwarderResource(fwd))
	}

	paginatedResponse(c, "forwarder", size, page, resources, count, err)
}

// CreateEVMChainRequest is a JSONAPI request for creating an EVM chain.
type CreateEVMForwarderRequest struct {
	EVMChainID *utils.Big     `json:"chainID"`
	Address    common.Address `json:"address"`
}

// Create adds a new EVM chain.
func (cc *EVMForwardersController) Create(c *gin.Context) {
	request := &CreateEVMForwarderRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	orm := forwarders.NewORM(cc.App.GetSqlxDB(), cc.App.GetLogger(), cc.App.GetConfig())
	fwd, err := orm.CreateForwarder(request.Address, *request.EVMChainID)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponseWithStatus(c, presenters.NewEVMForwarderResource(fwd), "forwarder", http.StatusCreated)
}

// Delete removes an EVM chain.
func (cc *EVMForwardersController) Delete(c *gin.Context) {
	id, err := stringutils.ToInt32(c.Param("fwdID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	orm := forwarders.NewORM(cc.App.GetSqlxDB(), cc.App.GetLogger(), cc.App.GetConfig())
	err = orm.DeleteForwarder(id)

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, nil, "forwarder", http.StatusNoContent)
}
