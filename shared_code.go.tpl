// attrReadMissingDiag represents diagnostic message on an attribute missing in the source object
type attrReadMissingDiag struct {
    Path string
}

func (d attrReadMissingDiag) Severity() diag.Severity {
    return diag.SeverityError
}

func (d attrReadMissingDiag) Summary() string {
    return "Error reading from Terraform object"
}

func (d attrReadMissingDiag) Detail() string {
    return fmt.Sprintf("A value for %v is missing in the source Terraform object Attrs", d.Path)
}

func (d attrReadMissingDiag) Equal(o diag.Diagnostic) bool {
    return (d.Severity() == o.Severity()) && (d.Summary() == o.Summary()) && (d.Detail() == o.Detail())
}

// attrReadConversionFailureDiag represents diagnostic message on a failed type conversion on read
type attrReadConversionFailureDiag struct {
    Path string
    Type string
}

func (d attrReadConversionFailureDiag) Severity() diag.Severity {
    return diag.SeverityError
}

func (d attrReadConversionFailureDiag) Summary() string {
    return "Error reading from Terraform object"
}

func (d attrReadConversionFailureDiag) Detail() string {
    return fmt.Sprintf("A value for %v can not be converted to %v", d.Path, d.Type)
}

func (d attrReadConversionFailureDiag) Equal(o diag.Diagnostic) bool {
    return (d.Severity() == o.Severity()) && (d.Summary() == o.Summary()) && (d.Detail() == o.Detail())
}

// attrWriteMissingDiag represents diagnostic message on an attribute missing in the target object
type attrWriteMissingDiag struct {
    Path string
}

func (d attrWriteMissingDiag) Severity() diag.Severity {
    return diag.SeverityError
}

func (d attrWriteMissingDiag) Summary() string {
    return "Error writing to Terraform object"
}

func (d attrWriteMissingDiag) Detail() string {
    return fmt.Sprintf("A value for %v is missing in the source Terraform object AttrTypes", d.Path)
}

func (d attrWriteMissingDiag) Equal(o diag.Diagnostic) bool {
    return (d.Severity() == o.Severity()) && (d.Summary() == o.Summary()) && (d.Detail() == o.Detail())
}

// attrWriteConversionFailureDiag represents diagnostic message on a failed type conversion on write
type attrWriteConversionFailureDiag struct {
    Path string
    Type string
}

func (d attrWriteConversionFailureDiag) Severity() diag.Severity {
    return diag.SeverityError
}

func (d attrWriteConversionFailureDiag) Summary() string {
    return "Error writing to Terraform object"
}

func (d attrWriteConversionFailureDiag) Detail() string {
    return fmt.Sprintf("A value for %v can not be converted to %v", d.Path, d.Type)
}

func (d attrWriteConversionFailureDiag) Equal(o diag.Diagnostic) bool {
    return (d.Severity() == o.Severity()) && (d.Summary() == o.Summary()) && (d.Detail() == o.Detail())
}

// attrWriteGeneralError represents diagnostic message on a generic error on write
type attrWriteGeneralError struct {
    Path string
    Err error
}

func (d attrWriteGeneralError) Severity() diag.Severity {
    return diag.SeverityError
}

func (d attrWriteGeneralError) Summary() string {
    return "Error writing to Terraform object"
}

func (d attrWriteGeneralError) Detail() string {
    return fmt.Sprintf("%s: %s", d.Path, d.Err.Error())
}

func (d attrWriteGeneralError) Equal(o diag.Diagnostic) bool {
    return (d.Severity() == o.Severity()) && (d.Summary() == o.Summary()) && (d.Detail() == o.Detail())
}