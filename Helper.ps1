# =========================================================================================
# Scaffold Operator
# =========================================================================================

kubebuilder init --domain kcloutie.com --repo github.com/kcloutie/tkn-github-status --owner "Ken Cloutier"  #--plugins=go/v4-alpha
kubebuilder create api --group status --version v1 --kind PipelineRunStatus

# ================================================================================================
# Update the API code
# ================================================================================================
make manifests

# ================================================================================================
# Install the CRD into the cluster
# ================================================================================================
make install